package main

import (
	"log"

	"github.com/aserto-dev/aserto-idp-plugin-okta/pkg/srv"
	"github.com/aserto-dev/idp-plugin-sdk/plugin"
)

func main() {

	oktaPlugin := srv.NewOktaPlugin()

	options := &plugin.PluginOptions{
		PluginHandler: oktaPlugin,
	}

	err := plugin.Serve(options)
	if err != nil {
		log.Println(err.Error())
	}
}
