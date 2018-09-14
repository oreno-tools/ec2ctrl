package main

import (
    "os"
    "fmt"
    "flag"
    "strings"
    "bytes"
    "encoding/csv"
    "encoding/json"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/service/ec2"

    "github.com/olekukonko/tablewriter"
)

const AppVersion = "0.0.6"

var (
    argProfile = flag.String("profile", "", "Profile 名を指定.")
    argRegion = flag.String("region", "ap-northeast-1", "Region 名を指定.")
    argEndpoint = flag.String("endpoint", "", "AWS API のエンドポイントを指定.")
    argInstances = flag.String("instances", "", "Instance ID 又は Instance Tag 名を指定.")
    argTags = flag.String("tags", "", "Tag Key 及び Tag Value を指定.")
    argStart = flag.Bool("start", false, "Instance を起動.")
    argStop = flag.Bool("stop", false, "Instance を停止.")
    argState = flag.Bool("state", false, "Instance の状態を出力.")
    argVersion = flag.Bool("version", false, "バージョンを出力.")
    argCsv = flag.Bool("csv", false, "CSV 形式で出力する")
    argJson = flag.Bool("json", false, "JSON 形式で出力する")
)

type Results struct {
    Instances      []Instance    `json:"instances"`
}

type Instance struct {
    Name           string `json:"name"`
    InstanceId     string `json:"instance_id"`
    InstanceType   string `json:"instance_type"`
    AZ             string `json:"az"`
    PrivateIp      string `json:"private_ip"`
    PublicIp       string `json:"public_ip"`
    Satate         string `json:"state"`
    InstanceStatus string `json:"instance_state"`
    SystemStatus   string `json:"system_state"`
}

func outputTbl(data [][]string) {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"tag:Name", "InstanceId", "InstanceType",
                             "AZ", "PrivateIp", "PublicIp", "State",
                             "InstanceStatus", "SystemStatus"})
    for _, value := range data {
        table.Append(value)
    }
    table.Render()
}

func outputCsv(data [][]string) {
    buf := new(bytes.Buffer)
    w := csv.NewWriter(buf)
    for _, record := range data {
        if err := w.Write(record); err != nil {
            fmt.Println("Write error: ", err)
            return
        }
        w.Flush()
    }
    fmt.Println(buf.String())
}

func outputJson(data [][]string) {
    var rs []Instance
    for _, record := range data {
        r := Instance{Name:record[0], InstanceId:record[1], InstanceType:record[2],
                      AZ:record[3], PrivateIp:record[4], PublicIp:record[5],
                      Satate:record[6], InstanceStatus:record[7], SystemStatus:record[8]}
        rs = append(rs, r)
    }
    rj := Results{
        Instances: rs,
    }
    b, err := json.Marshal(rj)
    if err != nil {
        fmt.Println("JSON Marshal error:", err)
        return
    }
    os.Stdout.Write(b)
}

func awsEc2Client(profile string, region string) *ec2.EC2 {
    var config aws.Config
    if profile != "" {
        creds := credentials.NewSharedCredentials("", profile)
        config = aws.Config{Region: aws.String(region), Credentials: creds, Endpoint: aws.String(*argEndpoint)}
    } else {
        config = aws.Config{Region: aws.String(region), Endpoint: aws.String(*argEndpoint)}
    }
    sess := session.New(&config)
    ec2Client := ec2.New(sess)
    return ec2Client
}

func getInstanceStatus(ec2Client *ec2.EC2, instanceId string) (string, string) {
    params := &ec2.DescribeInstanceStatusInput {
        InstanceIds: []*string {
            aws.String(instanceId),
        },
    }
    res, err := ec2Client.DescribeInstanceStatus(params)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    var instance_status string
    var system_status string
    if len(res.InstanceStatuses) > 0 {
        instance_status = *res.InstanceStatuses[0].InstanceStatus.Status
        system_status = *res.InstanceStatuses[0].SystemStatus.Status
    } else {
        instance_status = "N/A"
        system_status = "N/A"
    }

    return instance_status, system_status
}

func listInstances(ec2Client *ec2.EC2, instances []*string) {
    params := &ec2.DescribeInstancesInput {
        InstanceIds: instances,
    }
    res, err := ec2Client.DescribeInstances(params)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    allInstances := [][]string{}
    for _, r := range res.Reservations {
        for _, i := range r.Instances {
            var tag_name string
            for _, t := range i.Tags {
                if *t.Key == "Name" {
                    tag_name = *t.Value
                }
            }
            if i.PublicIpAddress == nil {
                i.PublicIpAddress = aws.String("Not assignment")
            }
            if i.PrivateIpAddress == nil {
                i.PrivateIpAddress = aws.String("Not assignment")
            }
            instance_status, system_status := getInstanceStatus(ec2Client, *i.InstanceId)
            instance := []string{
                tag_name,
                *i.InstanceId,
                *i.InstanceType,
                *i.Placement.AvailabilityZone,
                *i.PrivateIpAddress,
                *i.PublicIpAddress,
                *i.State.Name,
                instance_status,
                system_status,
            }
            allInstances = append(allInstances, instance)
        }
    }

    if *argCsv == true {
        outputCsv(allInstances)
    } else if *argJson == true {
        outputJson(allInstances)
    } else {
        outputTbl(allInstances)
    }
}

func startInstances(ec2Client *ec2.EC2, instances []*string) {
    params := &ec2.StartInstancesInput{
        InstanceIds: instances,
    }
    result, err := ec2Client.StartInstances(params)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    for _, r := range result.StartingInstances {
        fmt.Printf("%s を起動しました.\n", *r.InstanceId)
    }
}

func stopInstances(ec2Client *ec2.EC2, instances []*string) {
    params := &ec2.StopInstancesInput{
        InstanceIds: instances,
    }
    // fmt.Println(params)
    result, err := ec2Client.StopInstances(params)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    for _, r := range result.StoppingInstances {
        fmt.Printf("%s を停止しました.\n", *r.InstanceId)
    }
}

func stateInstances(ec2Client *ec2.EC2, instances []*string) {
    listInstances(ec2Client, instances)
}

func ctrlInstances(ec2Client *ec2.EC2, instances []*string, operation string) {
    listInstances(ec2Client, instances)

    switch operation {
    case "start":
        fmt.Print("上記のインスタンスを起動しますか?(y/n): ")
    case "stop":
        fmt.Print("上記のインスタンスを停止しますか?(y/n): ")
    }

    var stdin string
    fmt.Scan(&stdin)
    switch stdin {
    case "y", "Y":
        switch operation {
        case "start":
            fmt.Println("EC2 を起動します.")
            startInstances(ec2Client, instances)
        case "stop":
            fmt.Println("EC2 を停止します.")
            stopInstances(ec2Client, instances)
        }
    case "n", "N":
        fmt.Println("処理を停止します.")
        os.Exit(0)
    default:
        fmt.Println("処理を停止します.")
        os.Exit(0)
    }
}

func getInstanceIds(ec2Client *ec2.EC2, instances string) []*string {
    splitedInstances := strings.Split(instances, ",")
    res, err := ec2Client.DescribeInstances(nil)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    var instanceIds []*string
    for _, s := range splitedInstances {
        for _, r := range res.Reservations {
            for _, i := range r.Instances {
                for _, t := range i.Tags {
                    if *t.Key == "Name" {
                        if *t.Value == s {
                             instanceIds = append(instanceIds, aws.String(*i.InstanceId))
                        }
                    }
                }
                if *i.InstanceId == s {
                    instanceIds = append(instanceIds, aws.String(*i.InstanceId))
                }
            }
        }
    }
    return instanceIds
}

func main() {
    flag.Parse()

    if *argVersion {
      fmt.Println(AppVersion)
      os.Exit(0)
    }

    ec2Client := awsEc2Client(*argProfile, *argRegion)
    var instances []*string
    if *argInstances != "" {
        instances = getInstanceIds(ec2Client, *argInstances)
        if *argStart {
            ctrlInstances(ec2Client, instances, "start")
        } else if *argStop {
            ctrlInstances(ec2Client, instances, "stop")
        } else if *argState {
            stateInstances(ec2Client, instances)
        } else {
            fmt.Println("`-start`, `-stop`, `-state` 何れかのオプションを指定して下さい.")
            os.Exit(1)
        }
    } else {
        listInstances(ec2Client, nil)
    }
}
