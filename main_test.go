package main

import "testing"

//func TestGenerate(t *testing.T) {
//	op := &SingleVersionGenerator{
//		SingleVersionOptions: SingleVersionOptions{
//			//InputPackage: "sigs.k8s.io/controller-tools/testData/pkg/apis/fun",
//			InputPackage: "sigs.k8s.io/controller-tools/testData/pkg/apis/fun/v1alpha1",
//			Types:        []string{"Toy"},
//			//Flatten:      true,
//		},
//		WriterOptions: WriterOptions{
//			OutputPath: "../../../sigs.k8s.io/controller-tools/testData/config/crds/fun_v1alpha1_toy.yaml",
//			//OutputPath: "output.json",
//			OutputFormat: "yaml",
//		},
//		outputCRD: true,
//	}
//	op.Generate()
//}

func TestMultiVersionGenerate(t *testing.T) {
	op := &MultiVersionGenerator{
		MultiVersionOptions: MultiVersionOptions{
			//InputPackage: "sigs.k8s.io/controller-tools/testData/pkg/apis/fun",
			InputPackage: "sigs.k8s.io/controller-tools/testData/pkg/apis/fun",
			Types:        []string{"Toy"},
			//Flatten:      true,
		},
		WriterOptions: WriterOptions{
			OutputPath: "../../../sigs.k8s.io/controller-tools/testData/config/crds/fun_v1alpha1_toy.yaml",
			//OutputPath: "output.json",
			OutputFormat: "yaml",
		},
	}
	op.Generate()
}
