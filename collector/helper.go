package collector
 
import "strconv"
 
func toFloat(s string) float64 {
        // ubah string ke float
        f, _ := strconv.ParseFloat(s, 64)
        return f
}