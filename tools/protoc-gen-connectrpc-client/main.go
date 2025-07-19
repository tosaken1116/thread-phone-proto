package main

import (
	"flag"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	var flags flag.FlagSet

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			if err := generateFile(gen, f); err != nil {
				return err
			}
		}
		return nil
	})
}

func generateFile(gen *protogen.Plugin, file *protogen.File) error {
	for _, service := range file.Services {
		serviceName := string(service.Desc.Name())
		pkg := string(file.Desc.Package())
		pkgParts := strings.Split(pkg, ".")
		
		if len(pkgParts) < 2 {
			continue
		}
		
		domain := pkgParts[1]
		filename := "string_phone/" + domain + "/index.ts"
		
		g := gen.NewGeneratedFile(filename, file.GoImportPath)


		importPath := getImportPath(domain, serviceName)
		
		g.P("import { ", serviceName, " } from \"@/generated/", importPath, "\";")
		g.P("import { createClient } from \"@connectrpc/connect\";")
		g.P("import { createConnectTransport } from \"@connectrpc/connect-web\";")
		g.P("")
		g.P("export const ", serviceName, "ClientCreator = (baseUrl: string) => createClient(", serviceName, ", createConnectTransport({")
		g.P("  baseUrl,")
		g.P("}));")
	}

	return nil
}

func getPackagePath(file *protogen.File) string {
	pkg := string(file.Desc.Package())
	return strings.ReplaceAll(pkg, ".", "/") + "_connect"
}

func getImportPath(domain, serviceName string) string {
	switch serviceName {
	case "UserService":
		return domain + "/user_pb"
	case "StringItemService":
		if domain == "cards" {
			return domain + "/card_pb"
		}
		return domain + "/string_pb"
	}
	
	serviceFile := strings.ToLower(serviceName)
	if strings.HasSuffix(serviceFile, "service") {
		serviceFile = serviceFile[:len(serviceFile)-7]
	}
	return domain + "/" + serviceFile + "_pb"
}