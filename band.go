package learnermodel

import "strconv"

// BandByGrade gives the start band of a module for a learner's grade, from the bands of the
// version the attempt ran on, in their order (docs/architecture.md, sections 3.1 and 3.5): the
// band whose grades name the grade; for a grade no band names, the band nearest by the distance to
// its nearest grade, and the lower of two as near; without a grade, or in a module without grades,
// the first band. api puts it in the evidence it builds, and Append keeps it only for a module with
// no state yet.
//
// It refuses, with ErrInvalid, no bands, bands out of their places 1, 2, … and a band ID twice, as
// NewGraph does.
func BandByGrade(grade Grade, bands []Band) (BandRef, error) {
	if len(bands) == 0 {
		return BandRef{}, invalid("no bands")
	}
	ids := make([]string, 0, len(bands))
	for i, b := range bands {
		if int(b.Place) != i+1 {
			return BandRef{}, invalid("band " + b.ID + " at place " + strconv.Itoa(int(b.Place)) + ", not " + strconv.Itoa(i+1))
		}
		ids = append(ids, b.ID)
	}
	if err := idList("band", ids); err != nil {
		return BandRef{}, err
	}
	best, bestDistance := bands[0], int64(-1)
	if grade.Known {
		for _, b := range bands {
			for _, g := range b.Grades {
				d := int64(grade.Value) - int64(g)
				if d < 0 {
					d = -d
				}
				if bestDistance < 0 || d < bestDistance {
					best, bestDistance = b, d
				}
			}
		}
	}
	return BandRef{Place: best.Place, ID: best.ID}, nil
}
