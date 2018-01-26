#!/usr/bin/python2.7

from bs4 import BeautifulSoup
import urllib2
import MySQLdb
import time

def spider(url):
    request = urllib2.Request(url)
    request.add_header('User-Agent', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36')
    opener = urllib2.build_opener()
    #opener.addheaders.append(('Cookie',''))
    opener.addheaders.append(('Cookie','__cfduid=dff3d50d7a7fba6a53dd29b711e7d93521516259067; ASP.NET_SessionId=gde0u3rf5lxk53qfjmokhfjb; _ga=GA1.2.1085595957.1516259070; _gid=GA1.2.33895899.1516616011; __cflb=4215917907; cf_clearance=c1c1598e3452b4f9f37df8cd49e0e9cdac7b4f17-1516671681-10800'))
    f = opener.open(request)
    soup = BeautifulSoup(f.read().decode('utf-8'), "html5lib")
    return soup

def getAccount(db, cursor, addr):
    sql = "SELECT * FROM Accounts WHERE Address='%s'" % addr
    cursor.execute(sql)
    results = cursor.fetchall()
    if len(results)==0:
        fromSoup = spider("https://etherscan.io/token/EOS?a=%s" % addr)
        div = fromSoup.find(id='ContentPlaceHolder1_tr_tokenValue')
        if div==None:
            return 1
        td = div.findAll('td')
        if td==None or len(td)!=2:
            return 1
        value = td[1].text
        if value==None or value.strip()=='':
            return 1
        value = value.split(' ')[0]
        value = value.replace('$', '')
        value = value.replace(',', '')
        value = float(value)
        print addr, value
        
        name = fromSoup.find(id='ContentPlaceHolder1_divSummary')
        if name!=None:
            name = name.find('font')
            if name!=None:
                name = name.text
                if name!=None and name.strip()!='':
                    name = name.replace('Filtered By Individual Token Holder', '')
                    name = name.strip()
        if name==None or name.strip()=='':
            name = ''
        isql = "INSERT INTO Accounts(`Address`, `Name`, `Balance`) VALUE('%s', '%s', '%.2f')" % (addr, name, value)
        cursor.execute(isql)
        db.commit()
        time.sleep(0.2)
    return 0
        

db = MySQLdb.connect("127.0.0.1","eth","eth","eth" )

while True:
  try:
    soup = spider('https://etherscan.io/token/generic-tokentxns2?contractAddress=0x86fa049857e0209aa7d9e616f7eb3b3b78ecfdb0&a=&mode=')
    table = soup.find(attrs={'class':'table'})
    db.ping(True)
    cursor = db.cursor()
    for row in table.findAll("tr"):
        cells = row.findAll("td")
        if len(cells) == 7:
            txhash = cells[0].find('a').text
            if txhash==None or txhash.strip()=='':
                continue
            age = cells[1].find('span').get('title')
            if age==None or age.strip()=='':
                continue
            fromc = cells[2].find('a').text
            if fromc==None or fromc.strip()=='':
                continue
            toc = cells[4].find('a').text
            if toc==None or toc.strip()=='':
                continue
            quantity = cells[5].text
            if quantity==None or quantity.strip()=='':
                continue
            quantity = quantity.replace(',', '')
            quantity = float(quantity)
            sql = "INSERT INTO Transactions(`TxHash`, `Age`, `FromAcc`, `ToAcc`, `Quantity`) VALUE('%s', '%s', '%s', '%s', '%.2f')" % (txhash, age, fromc, toc, quantity)
            db.ping(True)
            cursor = db.cursor()
            try:
                cursor.execute(sql)
            except:
                print 'already exists'
            db.commit()
            time.sleep(0.2)
            getAccount(db, cursor, fromc)
            getAccount(db, cursor, toc)
  except Exception, e:
    print e
db.close()
        
