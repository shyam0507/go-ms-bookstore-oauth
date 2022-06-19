package cassandra

import (
	"time"

	"github.com/aws/aws-sigv4-auth-cassandra-gocql-driver-plugin/sigv4"
	"github.com/gocql/gocql"
)

var session *gocql.Session

func init() {
	cluster := gocql.NewCluster("cassandra.us-east-1.amazonaws.com:9142")
	var auth sigv4.AwsAuthenticator = sigv4.NewAwsAuthenticator()
	auth.Region = "us-east-1"
	auth.AccessKeyId = "AKIA37HUGRTIUNMAVI6L"
	auth.SecretAccessKey = "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"

	cluster.Authenticator = auth
	cluster.ConnectTimeout = time.Second * 10
	cluster.Keyspace = "oauth"

	cluster.SslOpts = &gocql.SslOptions{
		CaPath: "./sf-class2-root.crt",
	}
	cluster.Consistency = gocql.LocalQuorum
	cluster.DisableInitialHostLookup = true

	var err error
	if session, err = cluster.CreateSession(); err != nil {
		panic(err)
	}
}

func GetSession() *gocql.Session {
	return session
}
