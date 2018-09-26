default:
	@echo Makefile for My Golang Application
	@echo
	@echo Usage: make [task] [ARGS]
	@echo
	@echo Tasks:
	@python -c 'from tasks import tasks; tasks()' < Makefile

depend: # 依存パッケージの導入
	@gom install

test: # test テストの実行
	@gom test -v

build: build # バイナリをビルドする
	@./build.sh ec2ctrl.go

release: release # バイナリをリリースする. 引数に `_VER=バージョン番号` を指定する.
	@ghr -u inokappa -r ec2ctrl v${_VER} ./pkg/
