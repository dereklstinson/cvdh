package cvdh



//Point returns a copy of a point position
type Point interface{
	Point()([] int)
}
//NewPoint creates a new point
func NewPoint(p []int)Point{
	np:=new(point)
	np.p=make([]int,len(p))
	copy(np.p,p)
	return np
}

//Point is a point in a plain.
type point struct {
	p []int
}
func newpoint(p Point)point{
	return point{
		p:p.Point(),
	}
}
func (p *point)Point()([] int){
	np:=make([]int,len(p.p))
	copy(np,p.p)
	return np
}
