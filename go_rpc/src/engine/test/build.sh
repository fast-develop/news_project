cur_dir=`pwd`
echo "cur_dicr:" $cur_dir
project_dir=`pwd`
root_dir=`cd ../../../;pwd`
project_name=test
echo "root_dir" $root_dir

export GOPATH=$root_dir
go build -o $project_dir/test $project_dir/test.go
go build -o $project_dir/test_http2 $project_dir/test_http2.go
go build -o $project_dir/test_press $project_dir/test_press.go
go build -o $project_dir/test_http2_press $project_dir/test_http2_press.go

