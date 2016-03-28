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

// User Public Key
const SK_USER string = "432398415306986194693973996870836079581453988813"

// CA Public Key
const SK_CA string = "49336018324808093534733548840411752485726058527829630668967480568854756416567496216294919051910148686186622706869702321664465094703247368646506821015290302480990450130280616929226917246255147063292301724297680683401258636182185599124131170077548450754294083728885075516985144944984920010138492897272069257160"

// System Paramters
const sys_p string = "168199388701209853920129085113302407023173962717160229197318545484823101018386724351964316301278642143567435810448472465887143222934545154943005714265124445244247988777471773193847131514083030740407543233616696550197643519458134465700691569680905568000063025830089599260400096259430726498683087138415465107499"
const sys_q string = "959452661475451209325433595634941112150003865821"
const sys_g string = "94389192776327398589845326980349814526433869093412782345430946059206568804005181600855825906142967271872548375877738949875812540433223444968461350789461385043775029963900638123183435133537262152973355498432995364505138912569755859623649866375135353179362670798771770711847430626954864269888988371113567502852"

func checkParamValidity(dss_p *big.Int, dss_q *big.Int, dss_g *big.Int) int {

	p_check := (dss_p.ProbablyPrime(PROB_PRIME_LIMIT) && (dss_p.BitLen() >= DSS_P_MIN_LEN) && (dss_p.BitLen() <= DSS_P_MAX_LEN) && (dss_p.BitLen()%DSS_P_STRENGTH == 0))
	if !p_check {
		return 1
	}

	q_check := ((dss_q.BitLen() == DSS_Q_LEN) && (dss_q.ProbablyPrime(PROB_PRIME_LIMIT)) &&
		(new(big.Int).Cmp(new(big.Int).Mod(new(big.Int).Sub(dss_p, big.NewInt(1)), dss_q))) == 0)
	if !q_check {
		return 1
	}

	g_check := (new(big.Int).Exp(dss_g, dss_q, dss_p).Cmp(big.NewInt(1)) == 0)
	if !g_check {
		return 1
	}

	return 0
}

func verifyCert(dss_r_str string, dss_s_str string, dss_h_str string, dss_g *big.Int, dss_p *big.Int) int {
	dss_r := big.NewInt(0)
	dss_r.SetString(dss_r_str, 10)

	dss_s := big.NewInt(0)
	dss_s.SetString(dss_s_str, 10)

	dss_h := big.NewInt(0)
	dss_h.SetString(dss_h_str, 10)

	sys_q_bigInt := big.NewInt(0)
	sys_q_bigInt.SetString(sys_q, 10)

	ca_public_key := big.NewInt(0)
	ca_public_key.SetString(SK_CA, 10)

	if (dss_r.Cmp(big.NewInt(0)) == 1) && (dss_r.Cmp(sys_q_bigInt) == -1) && (dss_r.Cmp(big.NewInt(0)) == 1) && (dss_s.Cmp(sys_q_bigInt) == -1) {
		dss_u := new(big.Int).Mod((new(big.Int).Mul(dss_h, new(big.Int).ModInverse(dss_s, sys_q_bigInt))), sys_q_bigInt)

		dss_v := new(big.Int).Mod(new(big.Int).Mul(dss_r, new(big.Int).ModInverse(dss_s, sys_q_bigInt)), sys_q_bigInt)

		dss_w := new(big.Int).Exp(dss_g, dss_u, dss_p)
		fmt.Println("dss_w: ", dss_w)
		dss_w = new(big.Int).Mul(dss_w, new(big.Int).Exp(ca_public_key, dss_v, dss_p))
		dss_w = new(big.Int).Mod(dss_w, dss_p)
		dss_w = new(big.Int).Mod(dss_w, sys_q_bigInt)
		fmt.Println("dss_w: ", dss_w)

		if dss_w.Cmp(dss_r) == 0 {
			return 0
		} else {
			fmt.Println("dss_r: ", dss_r)
			fmt.Println("sys_q_bigInt: ", sys_q_bigInt)
			fmt.Println("dss_s: ", dss_s)
			fmt.Println("dss_h: ", dss_h)
			fmt.Println("dss_u: ", dss_u)
			fmt.Println("dss_v: ", dss_v)
			fmt.Println("dss_w: ", dss_w)
			fmt.Println("dss_r: ", dss_r)
		}
	}

	return -1
}

func main() {
	dss_p := big.NewInt(0)
	dss_p.SetString(sys_p, 10)

	dss_q := big.NewInt(0)
	dss_q.SetString(sys_q, 10)

	dss_g := big.NewInt(0)
	dss_g.SetString(sys_g, 10)

	if checkParamValidity(dss_p, dss_q, dss_g) != 0 {
		fmt.Println("The selection of either p, q, or g doesn't meet the DSS requirement")
		fmt.Println("Exiting...")
		os.Exit(1)
	}

	connection, _ := net.Dial("tcp", "localhost:8081")
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the user identity: ")
	user_id, _ := reader.ReadString('\n')
	fmt.Fprintf(connection, user_id+"\n")

	fmt.Print("Public key OK? [", SK_USER, "]: Press Enter to accept.")
	reader.ReadString('\n')
	fmt.Fprintf(connection, SK_USER+"\n")

	// Receive stuff from goca-server
	dss_r_str, _ := bufio.NewReader(connection).ReadString('\n')
	dss_s_str, _ := bufio.NewReader(connection).ReadString('\n')
	user_public_key_str, _ := bufio.NewReader(connection).ReadString('\n')
	expDate_str, _ := bufio.NewReader(connection).ReadString('\n')
	dss_h_str, _ := bufio.NewReader(connection).ReadString('\n')

	if verifyCert(dss_r_str, dss_s_str, dss_h_str, dss_g, dss_p) == 0 {
		fmt.Println("DSS Certificate is Valid!")
		fmt.Printf("dss_r = %v", dss_r_str)
		fmt.Printf("dss_s = %v\n", dss_s_str)
		fmt.Printf("user_public_key_str = %v\n", user_public_key_str)
		fmt.Printf("expiry date = %v\n", expDate_str)
		fmt.Printf("dss_hash = %v\n", dss_h_str)
	} else {
		fmt.Println("DSS Certificate is invalid!")
		fmt.Println("dss_r_str: ", dss_r_str)
		fmt.Println("dss_s_str: ", dss_s_str)
		fmt.Println("user_public_key: ", SK_USER)
		fmt.Println("expDate_str: ", expDate_str)
		fmt.Println("dss_h_str: ", dss_h_str)
	}
}
