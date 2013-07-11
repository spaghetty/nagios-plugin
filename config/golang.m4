##
## AC macros for using go lang
##

AC_DEFUN([CHECK_GOLANG],[
	AC_ARG_VAR(GOLANG, [ Go compiler full path ])

	AC_PATH_PROG([GOLANG],[go],[no])
	if test x$GOLANG != xno; then
	   found_golang="yes";
	fi
	AS_VAR_SET_IF([found_golang],
	[],
	[AC_MSG_ERROR([Cannot find golang installed. Install go and set GOROOT var])]
	)
])
