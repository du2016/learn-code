#  编译归档文件

cd x_src && go build -o x.a -i

# 拷贝到归档目录

mv x.a $GOPATH/pkg/darwin_amd64/github.com/du2016/learn-code/binary-pkg/

# 运行

go run main.go