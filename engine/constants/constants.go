package constants

const (
	// Window
	WindowWidth  = 1080
	WindowHeight = 720
)

const (
	// Algorithm names
	NoAlgo = "No Algorithm !TODO!"

	BVH = "Bounding Volume Hierarchy"
	SaP = "Sweep and Prune"
	GJK = "Gilbert-Johnson-Keerthi"
	SAT = "Separating Axis Theorem"
	LCP = "Linear Complementarity Problem"
	PGS = "Projected Gauss-Seidel"
	TGS = "Temporal Gauss-Seidel"
	MSI = "Method of sequential impulses"
	EPA = "Expanding Polytope Algorithm"
)

const (
	N   = " (No Parallel)"
	PT  = " (Parallel trivial)"
	PNT = " (Parallel non trivial)"

	ParallelPipeline   = "Parallel Pipeline"
	SequentialPipeline = "Sequential Pipeline"
)

var (
	AlgoType          = N
	SecondaryAlgoType = N
	ResolveAlgoType   = N
	Pipeline          = SequentialPipeline
)
