# 巡查系统

#### 分为两个模块
coco_server 作为连接所有集群nat机的跳板服务器

patrol_server 用于和用户交互，收集数据，分发脚本等功能
也是整个API的接口模块

###### 需要注意的是，patrol_server模块的Secret文件中需要加入私有的微信信息