#! /bin/bash
if [ $# -eq 0 ];then
	echo "启动容器: up"
	echo "关闭容器: down"
fi

case $1 in
	up)

		echo "启动 orderer 组织节点..."
		docker-compose -f orderer.yaml up -d
		echo "启动 dairy 组织节点..."
		docker-compose -f dairy.yaml up -d
		echo "启动 machining 组织节点..."
		docker-compose -f machining.yaml up -d
		echo "启动 sale 组织节点..."
		docker-compose -f sale.yaml up -d

		echo "==========查看进程=========="
		docker-compose -f orderer.yaml ps
		docker-compose -f dairy.yaml ps
		docker-compose -f machining.yaml ps
		docker-compose -f sale.yaml ps
		;;
	down)
		echo "========= 关闭 =========="
		docker-compose -f orderer.yaml down -v
		docker-compose -f dairy.yaml down -v
		docker-compose -f machining.yaml down -v
		docker-compose -f sale.yaml down -v
		docker rm -f `docker ps -aq`
		docker rmi -f $(docker images | grep "dev-" | awk '{print $ 3}')
		;;
	*)
		echo "args error, do nothing..."
esac

