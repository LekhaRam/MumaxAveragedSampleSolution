//+build ignore

package main

import (
	"github.com/mumax/3/data"
	. "github.com/mumax/3/engine"
	"github.com/mumax/3/oommf"
	"math"
	"os"
)

const Mu0 = 4 * math.Pi * 1e-7

func main() {

	defer InitAndClose()()

	Nx := 100
	Ny := 80
	Nz := 1

	SetGridSize(Nx, Ny, Nz)
	SetCellSize(1e-9, 1e-9, 1e-9)
	EdgeSmooth = 8

	setgeom( cuboid(100e-9, 80e-9, 1e-9) )
	saveas(geom, "cuboid")

	Msat  = 127.8e3
	Aex   = 4e-12
	alpha = 0.5

	// Initial magnetisation - ground state
	m = Uniform(1,0,0)
	// Excitation
	k0:= 0.5e6 // > ksw
	f := 100e9
	t0 := 50e-12

	B_ext = vector(.145, 0.0, 0.0) // in T for f between 4-6 GHz

	mask := NewSlice(3, NX, NY, NZ)
	x0 := 15e-9

	// limiting excitation till 5um along x
	for i := 0; i < 30; i++{
		for j := 0; j < NY; j++{
			r := index2coord(i, j, 0)
			x := r.X()
			y := r.Y()
			excite := (sin(k0*(x-x0))/(k0*(x-x0)))
			mask.set(0, i, j, 0, 0.0)
			mask.set(1, i, j, 0, excite)
			mask.set(2, i, j, 0, 0.0)
		}
	}
	// 0.036 in T - 25% of Bias field --increased from 5%
	B_ext.add(mask,0.036*sin(2*pi*f*(t-t0))/(2*pi*f*(t-t0)))
	OutputFormat = OVF1_TEXT
	// tableAdd(B_ext)
	probeStart:=50
	probeEnd:=100
	filename:=""

	sf_y:=10
	sf_x:=1

	for i := 0; i < 1; i++ {
		run(1e-12)
		M_probe := newVectorMask(Ny/sf_y,probeEnd-probeStart,1)
		for j := probeStart; j < probeEnd; j = j + sf_x {
			for k := 0; k < Ny/sf_y; k++ {
				M_probe.setVector(k,j-probeStart,0,m.Getcell(j,k*sf_y,0))
			}
		}
		filename = sprintf("YSampled_%[1]d", i)
		saveAs(M_probe,filename)
	}

}
