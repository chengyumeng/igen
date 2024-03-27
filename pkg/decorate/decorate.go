package decorate

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chengyumeng/igen/pkg/decorate/model"
)

type Decorator struct {
	Config

	pkg               *model.Package
	destination       *os.File
	generator         *generator
	outputPackagePath string
}

type Config struct {
	Source        string
	Destination   string
	SelfPackage   string
	ExplainFile   string
	CopyrightFile string
	Imports       string
	Logger        string
	Prom          string
}

func Default(cfg Config) Interface {
	return &Decorator{
		Config:      cfg,
		destination: os.Stdout,
	}
}

func (d *Decorator) LoadPackage() {
	pkg, err := sourceMode(d.Source)
	if err != nil {
		log.Fatalf("Unable to create directory: %v", err)
	}

	d.pkg = pkg
}

func (d *Decorator) SetDestination() {
	if len(d.Destination) > 0 {
		if err := os.MkdirAll(filepath.Dir(d.Destination), os.ModePerm); err != nil {
			log.Fatalf("Unable to create directory: %v", err)
		}
		f, err := os.Create(d.Destination)
		if err != nil {
			log.Fatalf("Failed opening destination file: %v", err)
		}
		d.destination = f
	}
}

func (d *Decorator) GenerateOutputPath() {
	// outputPackagePath represents the fully qualified name of the package of
	// the generated code. Its purposes are to prevent the module from importing
	// itself and to prevent qualifying type names that come from its own
	// package (i.e. if there is a type called X then we want to print "X" not
	// "package.X" since "package" is this package). This can happen if the audit
	// is output into an already existing package.
	d.outputPackagePath = d.SelfPackage
	if d.outputPackagePath == "" && d.Destination != "" {
		dstPath, err := filepath.Abs(filepath.Dir(d.Destination))
		if err == nil {
			pkgPath, err := parsePackageImport(dstPath)
			if err == nil {
				d.outputPackagePath = pkgPath
			} else {
				log.Println("Unable to infer -self_package from destination file path:", err)
			}
		} else {
			log.Println("Unable to determine destination file path:", err)
		}
	}
}

func (d *Decorator) CreateGenerator() {
	d.generator = &generator{
		filename:    d.Source,
		destination: d.Destination,
		auditNames:  map[string]string{},
		imports:     strings.Split(d.Imports, ","),
		logger:      d.Logger,
		prom:        d.Prom,
	}

	if d.ExplainFile != "" {
		bt, err := ioutil.ReadFile(d.ExplainFile)
		if err != nil {
			log.Fatalf("Failed reading explain file: %v", err)
		}

		if err := json.Unmarshal(bt, &d.generator.functions); err != nil {
			log.Fatalf("Failed unmarshal explain file: %v", err)
		}
	}

	if d.CopyrightFile != "" {
		header, err := ioutil.ReadFile(d.CopyrightFile)
		if err != nil {
			log.Fatalf("Failed reading copyright file: %v", err)
		}

		d.generator.copyrightHeader = string(header)
	}
}

func (d *Decorator) Decorate() {
	d.LoadPackage()
	d.SetDestination()
	defer d.destination.Close()
	d.GenerateOutputPath()
	d.CreateGenerator()
	if err := d.generator.Generate(d.pkg, d.outputPackagePath); err != nil {
		log.Fatalf("Failed generating audit: %v", err)
	}
	if _, err := d.destination.Write(d.generator.Output()); err != nil {
		log.Fatalf("Failed writing to destination: %v", err)
	}
}
