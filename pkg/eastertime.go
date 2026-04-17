// Copyright (c) 2013 eastertime Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// eastertime provides functions to get midnight date time of easter day
// for Catholic and Orthodox on Gregorian Calendar
// web github/vjeantet/eastertime

package eastertime

import (
	"errors"
	"math"
	"time"
)

// CatholicByYear returns time.Time for midnight on Catholic Easter of a given
// year use Delambre and Butcher's and return time based on Gregorian calendar
//
// @param int year The year as a number greater than 0
// @return time.Time The easter date as a time.Time
// @return errors.Error Error if exists, else nil
func CatholicByYear(year int) (time.Time, error) {
	var a, b, c, d, e, r int

	if year < 0 {
		return time.Now(), errors.New("year have to be greater than 0")
	}

	a = year % 19
	if year >= 1583 {
		var f, g, h, i, k, l, m int
		b = year / 100
		c = year % 100
		d = b / 4
		e = b % 4
		f = (b + 8) / 25
		g = (b - f + 1) / 3
		h = (19*a + b - d - g + 15) % 30
		i = c / 4
		k = c % 4
		l = (32 + 2*e + 2*i - h - k) % 7
		m = (a + 11*h + 22*l) / 451
		r = 22 + h + l - 7*m
	} else {
		b = year % 7
		c = year % 4
		d = (19*a + 15) % 30
		e = (2*c + 4*b - d + 34) % 7
		r = 22 + d + e
	}

	return time.Date(year, time.March, r, 0, 0, 0, 0, time.Local), nil
}

// OrthodoxByYear returns time.Time for midnight on Orthodox Easter of a given
// year use Meeus Julian algorithm and return time based on Gregorian calendar
//
// @param int year The year as a number greater than 325
// @return time.Time The easter date as a time.Time
// @return errors.Error Error if exists, else nil
func OrthodoxByYear(year int) (time.Time, error) {

	if year < 326 {
		return time.Now(), errors.New("year have to be greater than 325")
	}

	var a, b, c, d, e int
	var month time.Month
	var day float64

	a = year % 4
	b = year % 7
	c = year % 19
	d = (19*c + 15) % 30
	e = (2*a + 4*b - d + 34) % 7
	month = time.Month((d + e + 114) / 31)
	day = math.Floor(float64((d+e+114)%31 + 1))
	day = day + 13

	return time.Date(year, month, int(day), 0, 0, 0, 0, time.Local), nil
}
