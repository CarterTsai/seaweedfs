package filer2

import (
	"os"

	"github.com/spf13/viper"
	"github.com/chrislusf/seaweedfs/weed/glog"
)

const (
	FILER_TOML_EXAMPLE = `
# A sample TOML config file for SeaweedFS filer store

# local in memory, mostly for testing purpose
[memory]
enabled = false

[leveldb]
enabled = false
dir = "."					# directory to store level db files

[mysql]
enabled = true
server = "192.168.1.1"
port = 8080
username = ""
password = ""
database = ""
connection_max_idle = 100
connection_max_open = 100

[postgres]
enabled = false
server = "192.168.1.1"
port = 8080
username = ""
password = ""
database = ""
connection_max_idle = 100
connection_max_open = 100

`
)

var (
	Stores []FilerStore
)

func (f *Filer) LoadConfiguration() {

	// find a filer store
	viper.SetConfigName("filer")            // name of config file (without extension)
	viper.AddConfigPath(".")                // optionally look for config in the working directory
	viper.AddConfigPath("$HOME/.seaweedfs") // call multiple times to add many search paths
	viper.AddConfigPath("/etc/seaweedfs/")  // path to look for the config file in
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		glog.Fatalf("Failed to load filer.toml file from current directory, or $HOME/.seaweedfs/, or /etc/seaweedfs/" +
			"\n\nPlease follow this example and add a filer.toml file to " +
			"current directory, or $HOME/.seaweedfs/, or /etc/seaweedfs/:\n" + FILER_TOML_EXAMPLE)
	}

	glog.V(0).Infof("Reading filer configuration from %s", viper.ConfigFileUsed())
	for _, store := range Stores {
		if viper.GetBool(store.GetName() + ".enabled") {
			viperSub := viper.Sub(store.GetName())
			if err := store.Initialize(viperSub); err != nil {
				glog.Fatalf("Failed to initialize store for %s: %+v",
					store.GetName(), err)
			}
			f.SetStore(store)
			glog.V(0).Infof("Configure filer for %s from %s", store.GetName(), viper.ConfigFileUsed())
			return
		}
	}

	println()
	println("Supported filer stores are:")
	for _, store := range Stores {
		println("    " + store.GetName())
	}

	println()
	println("Please configure a supported filer store in", viper.ConfigFileUsed())
	println()

	os.Exit(-1)
}