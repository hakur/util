package alg

// YanghuiTriangle yanghui's traiangle , also known as pascal's triangle.  usage YanghuiTriangle(6, 1)
// YanghuiTriangle 杨辉三角形，也叫帕斯卡三角形. 使用 YanghuiTriangle(6, 1)
func YanghuiTriangle(height int, initNumber int) (geoData [][]int) {
	var lastLayerCount int

	for i := 0; i < height; i++ {
		layerData := make([]int, i+1)
		for j := 0; j <= i; j++ { // 每一层第一个数都是1，末端也是1
			if j > 0 {
				if j <= lastLayerCount-1 { // 没有超过超一层的数据个数
					layerData[j] = geoData[i-1][j-1] + geoData[i-1][j]
				} else {
					layerData[j] = geoData[i-1][lastLayerCount-1]
				}
			} else {
				layerData[j] = initNumber // 第一个数的值
			}
		}
		lastLayerCount = i + 1
		geoData = append(geoData, layerData)
	}
	return geoData
}
