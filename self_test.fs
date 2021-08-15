##
# Make sure that assertions are working
#  test assertTrue for direct case
#
func testOfAssertTrueDirect() {

    = assertTrue(true);
}

##
# Make sure that assertions are working
#  test assertTrue for reverse case
#
func testOfAssertTrueReverse() {

    try {

        assertTrue(false);

        return false;
    }
    catch($code, $msg) {

        return $msg == "java.lang.AssertionError: Expected logic true value";
    }
}

##
# Make sure that assertions are working
#  test assertFalse for direct case
#
func testOfAssertFalseDirect() {

    = assertFalse(false);
}

##
# Make sure that assertions are working
#  test assertFalse for reverse case
#
func testOfAssertFalseReverse() {

    try {

        assertFalse(true);

        return false;
    }
    catch($code, $msg) {

        return $msg == "java.lang.AssertionError: Expected logic false value";
    }
}

##
# Make sure that assertions are working
#  test assertEquals for numeric values
#
func testOfAssertEqualsIntegersDirect() {

    = assertEquals(1, 1);
}

##
# Make sure that assertions are working
#  test assertEquals for numeric values reverse case
#
func testOfAssertEqualsIntegersReverse() {

    try {

        assertEquals(1, 2);

        return false;
    }
    catch($code, $msg) {

        return $msg == "java.lang.AssertionError: Expected first argument equals to second: v1 = 1 v2 = 2";
    }
}

##
# Make sure that assertions are working
#  test assertEquals for numeric values
#
func testOfAssertEqualsFloatsDirect() {

    = assertEquals(1.123456, 1.123456);
}

##
# Make sure that assertions are working
#  test assertEquals for numeric values reverse case
#
func testOfAssertEqualsFloatsReverse() {

    try {

        assertEquals(1.123456, 1.12345);

        return false;
    }
    catch($code, $msg) {

        return $msg == "java.lang.AssertionError: Expected first argument equals to second: v1 = 1.123456 v2 = 1.12345";
    }
}

#@"when left operand is string then right should be cast to string to"
func test_3() {

    assertEquals("1", 1);

    return true;
}

#@"comparing to strings"
func test_4() {

    assertEquals("1", "1");

    return true;
}

#@"if left operand is a number then right will be casted to number"
func test_5() {

    assertEquals(1, "1");

    return true;
}

#@"make sure that true logic value equals to true logic value"
func test_6() {

    assertEquals(true, true);

    return true;
}

#@"... and false to false"
func test_7() {

    assertEquals(false, false);

    return true;
}

#@"comparing strings"
func test_8() {

    assertEquals("String", "String");

    return true;
}

#@"String comparision is case insensitive"
func test_9() {

    assertEquals("StRiNg", "sTrInG");

    return true;
}

#@"all the string should be in UTF8"
func test_10() {

    assertTrue("это текст UTF8" == "ЭТО ТЕКСТ utf8");

    return true;
}

#@"make sure that empty string is empty
#  (when explicit argument is missed func receives built-in NULL value reference)"
func test_11() {

    assertEquals("", );

    return true;
}

#@"implicit type casting true is casted to numerical 1 and string one is casted to 1"
func test_12() {

    assertTrue(true == "1");

    return true;
}

#@"implicit casting always tries to cast both operands to numbers (except situation when left operand is a string)"
func test_13() {

    assertTrue(true == "not empty");

    return true;
}

#@"when left operand is string"
func test_14() {

    assertTrue("true" == true);

    return true;
}

#@"casting false to numerical 0 but NULL is not a number (NaN)"
func test_15() {

    assertNotEquals(false, null);

    return true;
}

#@"casting false to 0"
func test_16() {

    assertTrue(false == 0);

    return true;
}

#@"casting true to 1"
func test_17() {

    assertTrue(true == 1);

    return true;
}

#@"testing operation priorities - first goes multiplication"
func test_18() {

    assertEquals(2 + 2 * 2, 6);

    return true;
}

#@"testing multiplication order - from left to right"
func test_19() {

    assertEquals(1 * 2 * 3, 6);

    return true;
}

#@"testing same priority operators order - from left to right"
func test_20() {

    assertEquals(2 * 2 / 2, 2);

    return true;
}

#@"...with more complicated expression"
func test_21() {

    assertEquals(2 * 2 + 2 / 2 - 1, 4);

    return true;
}

#@"testing sub expressions"
func test_22() {

    assertEquals(-2 * (2 + 2) * -2, 16);

    return true;
}

#@"... some more expressions"
func test_23() {

    assertEquals(2 * (-2 * 2) * 2, -16);

    return true;
}

#@"... and even that..."
func test_24() {

    assertEquals(2.0 / (2.0 * 2.0) * 1.0, .5);

    return true;
}

#@"testing operator < when it gives true"
func test_25() {

    assertTrue(1 < 2);

    return true;
}

#@"testing operator < when it gives false"
func test_26() {

    assertFalse(3 < 2);

    return true;
}

#@"testing operator > when it gives false"
func test_27() {

    assertFalse(1 > 2);

    return true;
}

#@"testing operator > when it gives true"
func test_28() {

    assertTrue(3 > 2);

    return true;
}

#@"testing operator >= when it gives true, case greater than"
func test_29() {

    assertTrue(3 >= 2);

    return true;
}

#@"testing operator >= when it gives true, case equals"
func test_30() {

    assertTrue(2 >= 2);

    return true;
}

#@"testing operator >= when it gives false"
func test_31() {

    assertFalse(1 >= 2);

    return true;
}

#@"testing operator <= when it gives true, case lesser than"
func test_32() {

    assertTrue(1 <= 2);

    return true;
}

#@"testing operator <= when it gives true, case equals"
func test_33() {

    assertTrue(4 <= 4);

    return true;
}

#@"testing operator <= when it gives false"
func test_34() {

    assertFalse(4 <= 3);

    return true;
}

#@"testing operator == when it gives true"
func test_35() {

    assertTrue(4 == 4);

    return true;
}

#@"testing operator == when it gives false"
func test_36() {

    assertFalse(4 == 3);

    return true;
}

#@"testing operator != when it gives true"
func test_37() {

    assertTrue(4 != 3);

    return true;
}

#@"testing operator != when it gives false"
func test_38() {

    assertFalse(4 != 4);

    return true;
}

#@"testing operator !(negation) not false is true"
func test_39() {

    assertTrue(!false);

    return true;
}

#@"... and not true is false"
func test_40() {

    assertFalse(!true);

    return true;
}

#@"true and true gives true"
func test_41() {

    assertTrue(true && true);

    return true;
}

#@"true and false gives false"
func test_42() {

    assertFalse(true && false);

    return true;
}

#@"false and true gives false"
func test_43() {

    assertFalse(false && true);

    return true;
}

#@"false and false gives false"
func test_44() {

    assertFalse(false && false);

    return true;
}

#@"false or false gives false"
func test_45() {

    assertFalse(false || false);

    return true;
}

#@"false or true gives true"
func test_46() {

    assertTrue(false || true);

    return true;
}

#@"true or false gives true"
func test_47() {

    assertTrue(true || false);

    return true;
}

#@"true or true gives true"
func test_48() {

    assertTrue(true || true);

    return true;
}

#@"testing ordering of logical operators true or false = true, true and true = true"
func test_49() {

    assertTrue(true || false && true);

    return true;
}

#@"testing ordering of logical operators false and false = false, false or true = true"
func test_50() {

    assertTrue(false && false || true);

    return true;
}

#@"testing priorities of logical operators - not false = true, true or false = true"
func test_51() {

    assertTrue(!false || false);

    return true;
}

#@"testing priorities of logical operators - not false = true, not false = true, true and true = true"
func test_52() {

    assertTrue(!false && ! false);

    return true;
}

#@"testing priorities of logical operators - 1 > 2 = false, 2 > 1 = true, false or true = true"
func test_53() {

    assertTrue(1 > 2 || 2 > 1);

    return true;
}

#@"testing priorities of logical operators - 1 > 2 = false, 2 > 1 = true, false and true = false"
func test_54() {

    assertFalse(1 > 2 && 2 > 1);

    return true;
}

#@"testing priorities of logical operators - 1 > 2 = false, 2 > 1 = true, not false = true, true and true = true"
func test_55() {

    assertTrue(!1 > 2 && 2 > 1);

    return true;
}

#@"testing priorities of logical operators"
func test_56() {

    assertTrue(!1 > 2 && !2 < 1);

    return true;
}

#@"testing priorities of logical operators"
func test_57() {

    assertTrue(!1 > 2);

    return true;
}

#@"testing priorities of logical operators"
func test_58() {

    assertFalse(!1 < 2);

    return true;
}

#@"make sure that it is logical"
func test_59() {

    assertEquals(2 < 3, true);

    return true;
}

#@"make sure that it is logical"
func test_60() {

    assertEquals(2 > 3, false);

    return true;
}

#@"testing if statement"
func test_61() {

    if(true) {

        assertTrue(true);

    }

    return true;
}

#@"testing if statement with return in block"
func test_61_1() {

    if(true) {

        assertTrue(true);

        return true;
    }

    return false;
}

#@"testing if statement with else"
func test_62() {

    if(false) {

        return false;
    }
    else {

        return true;
    }

}

func zero() {

    return 0;
}

#@"testing try/catch statement"
func test_63() {

    try {

        $n = 1 / zero();

        return false;
    }
    catch($code, $message) {

        return assertEquals($code, 3);
    }

}

#@"testing try/catch statement"
func test_64() {

    try {

        $n = 1 / zero();

        return false;
    }
    catch($code, $message) {

        = assertEquals($message, "Division by zero");
    }

}

#@"Test switch-case statement with general purpose"
func test_switch_case_1() {

    $variable = 1;

    switch($variable) {

        case(0) {
            return false;
        }

        case(1) {
            return true;
        }

        case(2) {
            return false;
        }
    }

}

#@"Test switch-case statement with reversed way"
func test_switch_case_2() {

    $a = 0;
    $b = 1;
    $c = 2;

    switch(1) {

        case($a) {
            return false;
        }

        case($b) {
            return true;
        }

        case($c) {
            return false;
        }
    }

    return false;
}

#@"Test switch-case statement with expressions inside case(...)"
func test_switch_case_3() {

    $a = 1;

    switch(2) {

        case($a * 3 + 1) {
            return false;
        }

        case(($a + 1) * 2 - 2) {
            return true;
        }

        case("xyz" + $a) {
            return false;
        }
    }

}

#@"Test switch-case statement with textual values"
func test_switch_case_4() {

    $a = "xyz";

    switch($a) {

        case("xyz" + $a) {
            return false;
        }

        case("xyz") {
            return true;
        }

        case($a > 1) {
            return false;
        }
    }

}

#@"Test switch-case with no returns"
func test_switch_case_5() {

    $a = "xyz";

    switch($a) {

        case("xyz" + $a) {
            return false;
        }

        case("xyz") {

            print(" ----> Hello world!");
        }

        case($a > 1) {
            return false;
        }
    }

    return true;
}

func sum($a, $b) {

    return $a + $b;
}

func test_function_1() {

    assertEquals(4, sum(2, 2));

    return true;
}

func test_function_2() {

    assertEquals(5, 1 + sum(2, 2));

    return true;
}

func test_function_3() {

    return assertEquals(5, sum(2, 2) + 1);
}

#@"No return statement function"
func test_function_4() {

    sum(1, 1);

    return true;
}

func test_function_5() {

    return assertEquals(7, 1 + sum(1 + 1, 2 + 1) + 1);
}

func test_strange_thing_1() {

    return assertEquals(true, true);
}

func test_while_statement() {

    $n = 10;

    while($n > 0) {

        print(" ----> $n = " + $n);

        $n = $n - 1;
    }

    return true;
}

func test_do_while_statement() {

    $n = 10;

    do {

        print(" ----> $n = " + $n);

        $n = $n - 1;
    } while ($n > 0);

    return true;
}

func test_pre_increment() {

    $a = 2;

    $b = ++ $a;

    assertEquals(3, $a);
    assertEquals(3, $b);

    return true;
}

func test_pre_decrement() {

    $a = 2;

    $b = -- $a;

    assertEquals(1, $a);
    assertEquals(1, $b);

    return true;
}

func test_post_increment() {

    $a = 2;

    $b = $a ++;

    assertEquals(3, $a);
    assertEquals(2, $b);

    return true;
}

func test_post_decrement() {

    $a = 2;

    $b = $a --;

    assertEquals(1, $a);
    assertEquals(2, $b);

    return true;
}

func test_for_1() {

    for($a = 0; $a < 5; $a ++) {

        print(" ----> $a = ", $a);
    }

    return assertEquals(5, $a);
}

func test_for_2() {

    for($a = 5; $a > -1; $a --) {

        print(" ----> $a = ", $a);
    }

    return assertEquals(-1, $a);
}

func test_comparison_numbers_1() {

    $a = 0.1 * 0.1 * 0.1;
    $b = 0.001;

    return assertTrue($a == $b);
}

func test_comparison_numbers_2() {

    $a = 0.1 * 0.1 * 0.1;
    $b = 0.001;

    return assertFalse($a < $b);
}

func test_comparison_numbers_3() {

    $a = 0.1 * 0.1 * 0.1;
    $b = 0.001;

    return assertFalse($a > $b);
}

func test_comparison_numbers_4() {

    $a = 0.1 * 0.1 * 0.1;
    $b = 0.001;

    return assertTrue($a >= $b);
}

func test_comparison_numbers_5() {

    $a = 0.1 * 0.1 * 0.1;
    $b = 0.001;

    return assertTrue($a <= $b);
}

func test_comparison_numbers_6() {

    $a = 1 + 0.1 * 0.1 * 0.1;
    $b = 1;

    return assertTrue($a >= $b);
}

func test_comparison_numbers_7() {

    $a = 0.1 * 0.1 * 0.1;
    $b = 0.01;

    return assertTrue($a <= $b);
}

func test_lambda() {

    $sum = func($a, $b) {

        return $a + $b;
    };

    print(" ----> Result: ", $sum(2, 4));

    return assertEquals(15, $sum(10, 5));
}

func func_returns_lambda($n) {

    //"=" is short for "return"
    return func($a, $c) {

        = $n + $a - $c;
    };
}

func test_lambda_result() {

    $sum = func_returns_lambda(10);

    print(" ----> Result: ", $sum(0, 5));

    = assertEquals(20, $sum(15, 5));
}

func test_func_reference() {

    $ref = @func_returns_lambda;

    $sum = $ref(20);

    print(" ----> Result: ", $sum(0, 5));

    = assertEquals(30, $sum(15, 5));
}

####
# decorator1 function is decorator definition
####
func decorator1($func) {

    return func($a, $b) {

        $res = $func($a, $b);

        $res = $res + 100;

        return $res;
    };
}

func decorator2($func) {

    return func($a, $b) {

        $res = $func($a, $b);

        $res += 200;

        return $res;
    };
}

#
# test conception without decoration
#
func test_pre_decorator() {

    $decorated = decorator1(func($a, $b){ = $a * $b; });

    return assertEquals(200, $decorated(10, 10));
}

@decorator1
@decorator2
func decoratedFunction($a, $b) {

    = $a + $b;
}

func test_decorated() {

    = assertEquals(330, decoratedFunction(40, -10));
}

func getEmpty($a, $b) {

    = $b;
}

func test_empty_expression() {

    = assertEquals(null, getEmpty(1,));
}

func test_parser_complicated_expression01($x) {

    $x = 1;

    = assertEquals(100 - ((10 - 10/$x) / 10), 100);
}

func test_parser_complicated_expression02() {

    $x = 10;

    = assertEquals(10 * $x + 20 * $x + 30 * $x + 40 * $x, 1000);
}

func test_muller_recurrence_roundoff_lack() {

    print(" ----> Muller's Recurrence - rounding test");

    $f = func($y, $z) {

        = 108 - ((815-1500/$z)/$y);
    };

    $x0 = 4;
    $x1 = 4.25;
    for($n = 0; $n < 20; $n ++) {

        print(" -----> Iteration: ", $n, " x0 = ", $x0);

        $tmp = $f($x1, $x0);

        $x0 = $x1;
        $x1 = $tmp;
    }

    print(" ----> Done: ", $n, " x0 = ", $x0);

    return $x0 < 5.0;
}

$test = "";

try {

    print("[i] Starting self test");

    $count = 0;

    while(($test = nextTest()) != false) {

        $count ++;

        print("[i] Test #", $count, ": ", $test, "():");

        if(execTest($test)) {

            print("[+] ", $test, "() passed");
        }
        else {

            print("[-] !!! ", $test, "() failed !!!");

            throw("Test failed!");
        }
    }

    print();
    print("[i] Total tests count: ", $count);
    print();

    = true;
}
catch($code, $msg) {


    print("[!] Test: ", $test);
    print("[!] Exception caught: ", $msg, "; code: ", $code);

    = false;
}