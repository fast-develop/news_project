cur_dir=`pwd`
echo "cur_dicr:" $cur_dir
project_dir=`pwd`
root_dir=`cd ../../;pwd`
project_name=engine
echo "root_dir" $root_dir

export GOPATH=$root_dir

protoc --go_out=plugins=grpc:$project_dir/proto/engine --proto_path=$project_dir/proto/engine $project_dir/proto/engine/engine.proto

find $project_dir/proto -name ".libs" | xargs rm -rf

mkdir -p $project_dir/_publish_dir/engine/log
cp -r $project_dir/release/* $project_dir/_publish_dir/engine
go build -gcflags '-N -l' -o $project_dir/_publish_dir/engine/bin/$project_name $project_dir/main/main.go
