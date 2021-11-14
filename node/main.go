package main
//
//func main() {
//	nodes := make([]*Node, 0)
//	n1 := NewNode("node1", "localhost", 8080, 1)
//	n2 := NewNode("node2", "localhost", 8081, 2)
//	//n3 := NewNode("node3", "localhost", 8082, 3)
//	//n4 := NewNode("node4", "localhost", 8083, 3)
//	//n5 := NewNode("node5", "localhost", 8084, 4)
//
//	var nodessss = n
//
//	nodes = append(nodes, n1)
//	nodes = append(nodes, n2)
//	//nodes = append(nodes, n3)
//	//nodes = append(nodes, n4)
//	//nodes = append(nodes, n5)
//
//	for _, node := range nodes {
//		go node.Init()
//
//		for _, no := range nodes {
//			if node.name == no.name {
//				continue
//			}
//			node.registerPeers(no.name, no.address, no.server.port)
//		}
//	}
//
//	for {
//
//	}
//}
