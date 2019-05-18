package config

import (
    "encoding/xml"
)

type HttpClientConfig struct {
    WeightAuto         int     `xml:"weight_auto,attr"`
    WeightMax          int     `xml:"weight_max,attr"`
    WeightAdjustPeriod int     `xml:"weight_adjust_period,attr"`
    WeightMaxOverload  float32 `xml:"weight_max_overload,attr"`
}

type KafkaClientConfig struct {
    TopicName   string `xml:"topicname,attr"`
    EnableSSL   string `xml:"enablessl,attr"`
    CertPath    string `xml:"certpath,attr"`
    ClusterNode int    `xml:"clusternode,attr"`
    Compression int8   `xml:"compression,attr"`
}

type DbClientConfig struct {
    DbName    string `xml:"dbname,attr"`
    TableName string `xml:"tablename,attr"`
    User      string `xml:"user,attr"`
    Password  string `xml:"password,attr"`
    Charset   string `xml:"charset,attr"`
}

type ClientInfo struct {
    Name          string   `xml:"name,attr"`
    Timeout       int      `xml:"timeout,attr"`
    TryNum        int      `xml:"trynum,attr"`
    MaxIdleConn   int      `xml:"maxidleconn,attr"`
    MaxActiveConn int      `xml:"maxactiveconn,attr"`
    MessageFormat string   `xml:"format,attr"`
    Addrs         []string `xml:"addr"`

    HttpClientConfig
    KafkaClientConfig
    DbClientConfig
}

type ClientConfig struct {
    XMLName      xml.Name      `xml:"clientconfig"`
    HttpClients  []*ClientInfo `xml:"httpclient"`
    RedisClients []*ClientInfo `xml:"redisclient"`
    MongoClients []*ClientInfo `xml:"mongoclient"`
    MysqlClients []*ClientInfo `xml:"mysqlclient"`
    KafkaClients []*ClientInfo `xml:"kafkaclient"`
    EtcdClients  []*ClientInfo `xml:"etcdclient"`
}

