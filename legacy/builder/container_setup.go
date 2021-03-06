/*
 * This file is part of Arduino Builder.
 *
 * Arduino Builder is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2015 Arduino LLC (http://www.arduino.cc/)
 */

package builder

import (
	bldr "github.com/arduino/arduino-cli/arduino/builder"
	"github.com/arduino/arduino-cli/legacy/builder/builder_utils"
	"github.com/arduino/arduino-cli/legacy/builder/i18n"
	"github.com/arduino/arduino-cli/legacy/builder/types"
	"github.com/arduino/go-paths-helper"
)

type ContainerSetupHardwareToolsLibsSketchAndProps struct{}

func (s *ContainerSetupHardwareToolsLibsSketchAndProps) Run(ctx *types.Context) error {
	commands := []types.Command{
		&AddAdditionalEntriesToContext{},
		&FailIfBuildPathEqualsSketchPath{},
		&HardwareLoader{},
		&PlatformKeysRewriteLoader{},
		&RewriteHardwareKeys{},
		&TargetBoardResolver{},
		&ToolsLoader{},
		&AddBuildBoardPropertyIfMissing{},
		&LibrariesLoader{},
	}

	ctx.Progress.Steps = ctx.Progress.Steps / float64(len(commands))

	for _, command := range commands {
		builder_utils.PrintProgressIfProgressEnabledAndMachineLogger(ctx)
		PrintRingNameIfDebug(ctx, command)
		err := command.Run(ctx)
		if err != nil {
			return i18n.WrapError(err)
		}
	}

	if ctx.SketchLocation != nil {
		// get abs path to sketch
		sketchLocation, err := ctx.SketchLocation.Abs()
		if err != nil {
			return i18n.WrapError(err)
		}

		// load sketch
		sketch, err := bldr.SketchLoad(sketchLocation.String(), ctx.BuildPath.String())
		if err != nil {
			return i18n.WrapError(err)
		}
		ctx.SketchLocation = paths.New(sketch.MainFile.Path)
		ctx.Sketch = types.SketchToLegacy(sketch)
	}

	commands = []types.Command{
		&SetupBuildProperties{},
		&LoadVIDPIDSpecificProperties{},
		&SetCustomBuildProperties{},
		&AddMissingBuildPropertiesFromParentPlatformTxtFiles{},
	}

	ctx.Progress.Steps = ctx.Progress.Steps / float64(len(commands))

	for _, command := range commands {
		builder_utils.PrintProgressIfProgressEnabledAndMachineLogger(ctx)
		PrintRingNameIfDebug(ctx, command)
		err := command.Run(ctx)
		if err != nil {
			return i18n.WrapError(err)
		}
	}

	return nil
}
