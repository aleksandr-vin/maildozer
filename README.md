# Maildozer -- an app to send nasty emails in a batch mode

## How to run

Write config file, for example `sample.yml`:
```
debug: yes
do-send: yes
mail-server: ginger.vinpn:25
from: alex.vinokurov@evecon.co
body-template: somefile.gtpl
subject-template: Evecon to {{.eventname}}, deer {{.shortname}}
to:
  -
    shortname: Bob
    fullname: Bob Tadam
    email: bob.tadam.test@evecon.co
    eventname: Megabolts'15 !! Yeah%$!
  -
    shortname: Dod
    fullname: Dod Tornton
    email: dod.tornton.test@evecon.co
    eventname: Supalulz'15 !! Yeah%$!
```
**NOTE: check the value of `do-send` parameter!**

Write email HTML template in Go `html/template` language, name it like `somefile.gtpl`:
```
<html>
  <body>
    <h1>Deer, {{.shortname}}!</h1>
  </body>
</html>
```

And call the app like `./maildozer sample.yml`, you should see smth. like this:
```
Maildozer is here to send your nasty mails, BOSS...
2015/09/10 01:39:53 Config data: map[body-template:somefile.gtpl subject-template:Evecon to {{.eventname}}, deer {{.shortname}} to:[map[shortname:Bob fullname:Bob Tadam email:bob.tadam.test@evecon.co eventname:Megabolts'15 !! Yeah%$!] map[fullname:Dod Tornton email:dod.tornton.test@evecon.co eventname:Supalulz'15 !! Yeah%$! shortname:Dod]] debug:true do-send:true mail-server:ginger.vinpn:25 from:alex.vinokurov@evecon.co]
2015/09/10 01:39:53 Mail server: ginger.vinpn:25
2015/09/10 01:39:53 From: alex.vinokurov@evecon.co
2015/09/10 01:39:53 Subject template: Evecon to {{.eventname}}, deer {{.shortname}}
2015/09/10 01:39:53 Body template file name: somefile.gtpl
2015/09/10 01:39:53 Do send emails: true
2015/09/10 01:39:54 Body: <html>
  <body>
    <h1>Deer, Bob!</h1>
  </body>
</html>

2015/09/10 01:39:54 Message 'Evecon to Megabolts'15 !! Yeah%$!, deer Bob' sent to Bob Tadam <bob.tadam.test@evecon.co>
2015/09/10 01:39:54 Body: <html>
  <body>
    <h1>Deer, Dod!</h1>
  </body>
</html>

2015/09/10 01:39:54 Message 'Evecon to Supalulz'15 !! Yeah%$!, deer Dod' sent to Dod Tornton <dod.tornton.test@evecon.co>
Accomplished

```
