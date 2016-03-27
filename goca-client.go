package main

import (
	"bufio"
	"fmt"
	"math/big"
	"net"
	"os"
)

const PROB_PRIME_LIMIT int = 50
const DSS_P_MIN_LEN int = 512
const DSS_P_MAX_LEN int = 1024
const DSS_P_STRENGTH int = 64
const DSS_Q_LEN int = 160

// System Paramters
const sys_p string = "168199388701209853920129085113302407023173962717160229197318545484823101018386724351964316301278642143567435810448472465887143222934545154943005714265124445244247988777471773193847131514083030740407543233616696550197643519458134465700691569680905568000063025830089599260400096259430726498683087138415465107499"
const sys_q string = "959452661475451209325433595634941112150003865821"
const sys_g string = "94389192776327398589845326980349814526433869093412782345430946059206568804005181600855825906142967271872548375877738949875812540433223444968461350789461385043775029963900638123183435133537262152973355498432995364505138912569755859623649866375135353179362670798771770711847430626954864269888988371113567502852"

func checkParamValidity(dss_p *big.Int, dss_q *big.Int, dss_g *big.Int) int {

	p_check := (dss_p.ProbablyPrime(PROB_PRIME_LIMIT) && (dss_p.BitLen() >= DSS_P_MIN_LEN) && (dss_p.BitLen() <= DSS_P_MAX_LEN) && (dss_p.BitLen()%DSS_P_STRENGTH == 0))
	if !p_check {
		return 1
	}

	// FIXME: This is the correct way to use the big.Int package
	//(*big.Int).Mod(big.NewInt(0), big.NewInt(1), big.NewInt(2))
	// %s/big.NewInt(0) or otherdss.p/q/g/ / (*big.Int)

	q_check := ((dss_q.BitLen() == DSS_Q_LEN) && (dss_q.ProbablyPrime(PROB_PRIME_LIMIT)) &&
		(big.NewInt(0).Cmp((big.NewInt(0).Mod(big.NewInt(0).Sub(dss_p, big.NewInt(1)), dss_q))) == 0))
	if !q_check {
		return 1
	}

	g_check := (dss_g.Exp(dss_g, dss_q, dss_p).Cmp(big.NewInt(1)) == 0)
	if !g_check {
		return 1
	}

	return 0
}

func main() {
	dss_p := *big.NewInt(0)
	dss_p.SetString(sys_p, 10)

	dss_q := *big.NewInt(0)
	dss_q.SetString(sys_q, 10)

	dss_g := *big.NewInt(0)
	dss_g.SetString(sys_g, 10)

	if checkParamValidity(&dss_p, &dss_q, &dss_g) != 0 {
		fmt.Println("The selection of either p, q, or g doesn't meet the DSS requirement")
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	connection, _ := net.Dial("tcp", "localhost:8081")
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the user identity: ")
	user_id, _ := reader.ReadString('\n')
	fmt.Fprintf(connection, user_id+"\n")

	fmt.Print("Enter the user's public key: ")
	user_public_key, _ := reader.ReadString('\n')
	fmt.Fprintf(connection, user_public_key+"\n")

	msg, _ := bufio.NewReader(connection).ReadString('\n')
	fmt.Print("Message from server: " + msg)
}
