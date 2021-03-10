
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUFFIXES = .go .c
OBJEXT= none

sbin_PROGRAMS = mqbot
mqbot_SOURCES = mqbot.go

GOFLAGS = -ldflags="-s -w"

EXTRA_mqbot_SOURCES = \
	mqbot.go

#EXTRA_mqbot_SOURCES += \
#	server/bundle/public.go

mqbot$(EXEEXT): $(mqbot_SOURCES) $(EXTRA_mqbot_SOURCES)
	$(GO) build $(GOFLAGS) -o mqbot$(EXEEXT) $(mqbot_SOURCES)

EXTRA_DIST = \
	README.md \
	go.mod \
	go.sum

#EXTRA_DIST += \
#	samples/mqbot.yml

clean-local:
	rm -rf autom4te.cache

install-data-local:
	test -z $(DESTDIR)$(APP_LIBDIR) || $(MKDIR_P) $(DESTDIR)$(APP_LIBDIR)
	test -z $(DESTDIR)$(APP_CONFDIR) || $(MKDIR_P) $(DESTDIR)$(APP_CONFDIR)
	test -z $(DESTDIR)$(APP_LOGDIR) || $(MKDIR_P) $(DESTDIR)$(APP_LOGDIR)
	test -z $(DESTDIR)$(APP_RUNDIR) || $(MKDIR_P) $(DESTDIR)$(APP_RUNDIR)
#	test -z $(DESTDIR)$(APP_DATADIR) || $(MKDIR_P) $(DESTDIR)$(APP_DATADIR)
#EOF
