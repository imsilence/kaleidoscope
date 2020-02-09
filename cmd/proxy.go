package cmd

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/kardianos/osext"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/imsilence/kaleidoscope/proxy"
	"github.com/imsilence/kaleidoscope/proxy/config"
)

var proxyCfgFile string

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "kaleidoscope proxy",
	Long:  `kaleidoscope proxy`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if proxyCfgFile != "" {
			viper.SetConfigFile(proxyCfgFile)
		} else {
			binHome, err := osext.ExecutableFolder()
			if err != nil {
				return err
			}
			viper.AddConfigPath(filepath.Join(binHome, "etc"))
			viper.AddConfigPath(filepath.Join(".", "etc"))
			viper.SetConfigName("proxy")
		}

		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{
			"file": viper.ConfigFileUsed(),
		}).Info("Using config file")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		proxyConfig := new(config.ProxyConfig)
		if err := viper.Unmarshal(proxyConfig); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("error proxy config unmarshal")
			return err
		}
		logrus.WithFields(logrus.Fields{
			"config": proxyConfig,
			"pid":    os.Getpid(),
		}).Debug("proxy starting...")
		proxy.DefaultManager.Init(proxyConfig)
		proxy.DefaultManager.Start()
		stop := make(chan os.Signal)
		signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT)
		<-stop
		proxy.DefaultManager.Stop()
		logrus.Debug("proxy stoping")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
	proxyCmd.Flags().StringVarP(&proxyCfgFile, "config", "c", "", "config file")

}
