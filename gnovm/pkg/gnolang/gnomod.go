package gnolang

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gnolang/gno/gnovm/pkg/gnomod"
	"github.com/gnolang/gno/gnovm/pkg/packages"
	"github.com/gnolang/gno/tm2/pkg/std"
)

const gnomodTemplate = `{{/*
This is a comment in a Go template in pkg/gnolang/gnomod.go.
The gnomodTemplate is used with the 'text/template' package
to generate the final gno.mod file. */}}
module {{.PkgPath}}

gno {{.GnoVersion}}`

func GenGnoModLatest(pkgPath string) string  { return genGnoMod(pkgPath, GnoVerLatest) }
func GenGnoModTesting(pkgPath string) string { return genGnoMod(pkgPath, GnoVerTesting) }
func GenGnoModDefault(pkgPath string) string { return genGnoMod(pkgPath, GnoVerDefault) }
func GenGnoModMissing(pkgPath string) string { return genGnoMod(pkgPath, GnoVerMissing) }

func genGnoMod(pkgPath string, gnoVersion string) string {
	buf := new(bytes.Buffer)
	tmpl := template.Must(template.New("").Parse(gnomodTemplate))
	err := tmpl.Execute(buf, struct {
		PkgPath    string
		GnoVersion string
	}{pkgPath, gnoVersion})
	if err != nil {
		panic(fmt.Errorf("generating gno.mod: %w", err))
	}
	return string(buf.Bytes())
}

const (
	GnoVerLatest  = `0.9` // current version
	GnoVerTesting = `0.9` // version of our tests
	GnoVerDefault = `0.9` // auto generated gno.mod
	GnoVerMissing = `0.0` // missing gno.mod, !autoGnoMod XXX
)

// ========================================
// Parses and checks the gno.mod file from mpkg.
// To generate default ones, GenGnoMod*().
//
// Results:
//   - mod: the gno.mod file, or nil if not found.
//   - err: wrapped error, or nil if file not found.
func ParseCheckGnoMod(mpkg *std.MemPackage) (mod *gnomod.File, err error) {
	if IsStdlib(mpkg.Path) {
		// stdlib/extern packages are assumed up to date.
		modstr := GenGnoModLatest(mpkg.Path)
		mod, _ = gnomod.ParseBytes("<stdlibs_autogenerated>/gno.mod", []byte(modstr))
	} else if mpkg.GetFile("gno.mod") == nil {
		// gno.mod doesn't exist.
		return nil, nil
	} else if mod, err = gnomod.ParseMemPackage(mpkg); err != nil {
		// error parsing gno.mod.
		err = fmt.Errorf("%s/gno.mod: parse error %w", mpkg.Path, err)
	} else if mod.Gno == nil {
		// gno.mod was never specified; set missing.
		mod.SetGno(GnoVerMissing)
	} else if mod.Gno.Version == GnoVerLatest {
		// current version, nothing to do.
	} else {
		panic("unsupported gno version " + mod.Gno.Version)
	}
	return
}

// ========================================
// ReadPkgListFromDir() lists all gno packages in the given dir directory.
func ReadPkgListFromDir(dir string) (gnomod.PkgList, error) {
	var pkgs []gnomod.Pkg

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		modPath := filepath.Join(path, "gno.mod")
		data, err := os.ReadFile(modPath)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}

		mod, err := gnomod.ParseBytes(modPath, data)
		if err != nil {
			return fmt.Errorf("parse: %w", err)
		}
		mod.Sanitize()
		if err := mod.Validate(); err != nil {
			return fmt.Errorf("failed to validate gno.mod in %s: %w", modPath, err)
		}

		pkg, err := ReadMemPackage(path, mod.Module.Mod.Path)
		if err != nil {
			// ignore package files on error
			pkg = &std.MemPackage{}
		}

		importsMap, err := packages.Imports(pkg, nil)
		if err != nil {
			// ignore imports on error
			importsMap = nil
		}
		importsRaw := importsMap.Merge(packages.FileKindPackageSource, packages.FileKindTest, packages.FileKindXTest)

		imports := make([]string, 0, len(importsRaw))
		for _, imp := range importsRaw {
			// remove self and standard libraries from imports
			if imp.PkgPath != mod.Module.Mod.Path &&
				!IsStdlib(imp.PkgPath) {
				imports = append(imports, imp.PkgPath)
			}
		}

		pkgs = append(pkgs, gnomod.Pkg{
			Dir:     path,
			Name:    mod.Module.Mod.Path,
			Draft:   mod.Draft,
			Imports: imports,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return pkgs, nil
}
