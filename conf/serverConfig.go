package conf

const (
	// 环境
	EnvRelease = "rls:"
	EnvDev     = "dev:"
	EnvTest    = "test:"

	// 模块配置文件在etcd中的key_name
	GatewayConfKey = "gw_ck"   // 网关在etcd中的key名称
	ServeConfKey   = "srv_ck"  // 服务(ws/sse)在etcd中的key名称
	AccountConfKey = "acc_ck"  // 用户模块在etcd中的key名称
	MsgConfKey     = "msg_ck"  // 消息模块在etcd中的key名称
	PlazaConfKey   = "plz_ck"  // 广场模块在etcd中的key名称
	CleanerConfKey = "cln_ck"  // 清理模块在etcd中的key名称
	FileConfKey    = "file_ck" // 文件模块在etcd中的key名称
	CMSConfKey     = "cms_ck"  // 管理后台在etcd中的key名称
)

// 网关模块配置
type GatewayConf struct {
	DBHost          string `json:"db_host"`           // sql型数据库ip
	DBPort          string `json:"db_port"`           // sql型数据库port
	DBUser          string `json:"db_user"`           // sql型数据库用户
	DBPwd           string `json:"db_pwd"`            // sql型数据库密码
	RedisClusterUrl string `json:"redis_cluster_url"` // redis集群地址
	RedisClusterPwd string `json:"redis_cluster_pwd"` // redis密码
	NatsUrl         string `json:"nats_url"`          // nats的url
	NatsPwd         string `json:"nats_pwd"`          // nats的密码
	ESUrl           string `json:"es_url"`            // Elasticsearch的url
	ESPwd           string `json:"es_pwd"`            // Elasticsearch的连接密码
}

// 服务模块的配置
type ServeConf struct {
	natsConf
	redisConf
}

// 用户模块配置
type AccountConf struct {
	sqlDb
	natsConf
	redisConf
	esConf
}

// 消息模块配置
type MsgConf struct {
	sqlDb
	natsConf
	redisConf
	esConf
}

// 广场模块配置
type PlazaConf struct {
	sqlDb
	natsConf
	redisConf
	esConf
}

// 清理模块配置
type CleanerConf struct {
	redisConf
	natsConf
}

// 文件模块配置
type FileConf struct {
	sqlDb
	redisConf
}

// 后台模块配置
type CMSConf struct {
	sqlDb
	redisConf
}

type redisConf struct {
	RedisNodes []string `json:"redis_nodes"`
	RedisUser  string   `json:"redis_user"`
	RedisPwd   string   `json:"redis_pwd"`
}

type natsConf struct {
	NatsUrl string `json:"nats_url"` // nats的url
	NatsPwd string `json:"nats_pwd"` // nats的密码
}

type sqlDb struct {
	DBHost string `json:"db_host"` // sql型数据库ip
	DBPort string `json:"db_port"` // sql型数据库port
	DBUser string `json:"db_user"` // sql型数据库用户
	DBPwd  string `json:"db_pwd"`  // sql型数据库密码
	DBName string `json:"db_name"` // sql型数据库的名称
}

type esConf struct {
	ESUrl  string `json:"es_url"`  // Elasticsearch的url
	ESUser string `json:"es_user"` // Elasticsearch的user
	ESPwd  string `json:"es_pwd"`  // Elasticsearch的连接密码
}
