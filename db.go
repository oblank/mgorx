package mgorx

import (
	"github.com/revel/config"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2"
	"os"
	"strings"
	"time"
)

var (
	mgoSession    *mgo.Session
	database_name string
)

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error

		//		url, found := revel.Config.String("db.url")
		//		if !found {
		//			panic("db.url not set in config")
		//		}
		//		database_name, found = revel.Config.String("db.name")
		//		if !found {
		//			panic("db.name not set in config")
		//		}
		//		mgoSession, err = mgo.Dial(url)
		//		if err != nil {
		//			panic(err) // no, not really
		//		}

		//判断是否是系统的分隔符
		separator := "/"
		if os.IsPathSeparator('\\') {
			separator = "\\"
		} else {
			separator = "/"
		}

		config_file := (revel.BasePath + "/conf/databases.conf")
		config_file = strings.Replace(config_file, "/", separator, -1)
		c, _ := config.ReadDefault(config_file)
		var section string
		if revel.DevMode {
			section = "database_dev"
		} else {
			section = "database_prod"
		}

		//TODO 连接池，mongosession copy, defer close
		//READ
		read_host, _ := c.String(section, "db.read.host")
		read_port, _ := c.String(section, "db.read.port")
		read_username, _ := c.String(section, "db.read.username")
		read_password, _ := c.String(section, "db.read.password")
		read_dbname, _ := c.String(section, "db.read.dbname")
		read_host_port := read_host + ":" + read_port
		mongoDBDialInfo4Read := &mgo.DialInfo{
			Addrs:    []string{read_host_port},
			Timeout:  60 * time.Second,
			Database: read_dbname,
			Username: read_username,
			Password: read_password,
		}
		// Create a session which maintains a pool of socket connections
		//		var err error
		mgoSession, err = mgo.DialWithInfo(mongoDBDialInfo4Read)
		if err != nil {
			revel.ERROR.Printf("DB_Read错误: %v", err)
			panic(err)
		}
		//		// Optional. Switch the session to a monotonic behavior.
		//		this.DBSessionRead.SetMode(mgo.Monotonic, true)
		//		this.DBRead = this.DBSessionRead.DB(read_dbname)
		//
		//		//WRITE
		//		wirte_host, _ := c.String(section, "db.write.host")
		//		write_port, _ := c.String(section, "db.write.port")
		//		write_username, _ := c.String(section, "db.write.username")
		//		write_password, _ := c.String(section, "db.write.password")
		//		write_dbname, _ := c.String(section, "db.write.dbname")
		//		wirte_host_port := wirte_host + ":" + write_port
		//		mongoDBDialInfo4Write := &mgo.DialInfo{
		//			Addrs:    []string{wirte_host_port},
		//			Timeout:  60 * time.Second,
		//			Database: write_dbname,
		//			Username: write_username,
		//			Password: write_password,
		//		}
		//		// Create a session which maintains a pool of socket connections
		//		this.DBSessionWrite, err = mgo.DialWithInfo(mongoDBDialInfo4Write)
		//		if err != nil {
		//			revel.ERROR.Printf("DB_Read错误: %v", err)
		//			panic(err)
		//		}
		//		// Optional. Switch the session to a monotonic behavior.
		//		this.DBSessionWrite.SetMode(mgo.Monotonic, true)
		//		this.DBWrite = this.DBSessionWrite.DB(write_dbname)

	}

	return mgoSession.Copy()
	//	return mgoSession.Clone()
}
