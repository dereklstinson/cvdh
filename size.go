package cvdh  

//Size returns a copy of the dim lengths.
type Size interface{
	Size()([] int)
//	Position()(x,y int)
}
//CreateSize creates a size
func CreateSize(s []int)Size{
ns:=new(size)
ns.s=make([]int,len(s))
copy(ns.s,s)
return ns
}

//SubSizes returns a subtraction of Subsizes of same dim.
//It assums you know what you are doing and doesn't check for error
func SubSizes(a,b Size)Size{
	na:=newsize(a)
nb:=newsize(b)

for i:=range na.s{
na.s[i]-=nb.s[i]
}
return na
}

type size struct {
	s	[]int
	}
	func newsize(s Size)*size{
		return &size{
			s:s.Size(),
		}
		
	}
	func (s *size)Size()([] int){
		ns:=make([]int,len(s.s))
		copy(ns,s.s)
		return ns
	}
	
	