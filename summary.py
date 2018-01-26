#!/usr/bin/python2.7

from bs4 import BeautifulSoup
import urllib2
import MySQLdb
import time
import datetime

def spider(url):
    request = urllib2.Request(url)
    request.add_header('User-Agent', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36')
    opener = urllib2.build_opener()
    opener.addheaders.append(('Cookie','__cfduid=dff3d50d7a7fba6a53dd29b711e7d93521516259067; ASP.NET_SessionId=gde0u3rf5lxk53qfjmokhfjb; _ga=GA1.2.1085595957.1516259070; _gid=GA1.2.33895899.1516616011; __cflb=4215917907; cf_clearance=c1c1598e3452b4f9f37df8cd49e0e9cdac7b4f17-1516671681-10800'))
    f = opener.open(request)
    soup = BeautifulSoup(f.read().decode('utf-8'), "html5lib")
    return soup

    
def getTransaction(addr):
    summary = {}
    transfer = {}
    
    flag = True
    page = 0
    while (flag):
        page = page + 1
        print page
        url = "https://etherscan.io/token/generic-tokentxns2?contractAddress=0x86fa049857e0209aa7d9e616f7eb3b3b78ecfdb0&mode=&a=%s&p=%d" % (addr, page)
        fromSoup = spider(url)
        div = fromSoup.find(id='maindiv')
        if div==None:
            return 1
        tb = div.find('table')
        tb = tb.find('tbody')
        tr = tb.findAll('tr')
        if tr==None or len(tr)==0:
            return 1
        for row in tr:
            cells = row.findAll('td')
            if len(cells)==7:
                age = cells[1].find('span').get('title')
                if age==None or age.strip()=='':
                    continue
                age = age.split(' ')[0]
                date_time = datetime.datetime.strptime(age,'%b-%d-%Y')
                day = date_time.strftime('%Y-%m-%d')
                if day<'2018-01-01':
                    flag = False
                    break
                
                type = cells[3].find('span').text
                if type==None or type.strip()=='':
                    continue
                    
                targetCell = cells[2]
                if "IN" in type:
                    type = "in"
                else:
                    type = "out"
                    targetCell = cells[4]
                target = targetCell.find('a').text
                if target==None or target.strip()=='':
                    continue
                quantity = cells[5].text
                if quantity==None or quantity.strip()=='':
                    continue
                quantity = quantity.replace(',', '')
                quantity = float(quantity)
                
                if not summary.has_key(day):
                    summary[day] = {}
                if not transfer.has_key(type):
                    transfer[type] = {}
                    
                if summary[day].has_key(type):
                    summary[day][type] = summary[day][type] + quantity
                else:
                    summary[day][type] = quantity
                    
                if transfer[type].has_key(target):
                    transfer[type][target] = transfer[type][target] + quantity
                else:
                    transfer[type][target] = quantity
            
        time.sleep(0.5)

    db = MySQLdb.connect("127.0.0.1","eth","eth","eth" )
    cursor = db.cursor()
    for d in summary:
        inAmount = 0
        outAmount = 0
        if "in" in summary[d]:
            inAmount = summary[d]["in"]
        if "out" in summary[d]:
            outAmount = summary[d]["out"]
        sql = "INSERT INTO Summary(`Account`,`Day`,`InAmount`,`OutAmount`) VALUE('%s','%s','%.2f','%.2f')" %(addr, d, inAmount, outAmount)
        print sql
        cursor.execute(sql)
        db.commit()
        
    print ""
        
    for type in transfer:
        for target in transfer[type]:
            sql = "INSERT INTO Transfer(`Account`,`Type`,`Target`,`Number`) VALUE('%s','%s','%s','%.2f')" %(addr, type, target, transfer[type][target])
            print sql
            cursor.execute(sql)
            db.commit()
        
    sql = "UPDATE Accounts set flag='y' where Address='%s'" %addr
    cursor.execute(sql)
    db.commit()
    db.close()
    return 0
        

while True:
    db = MySQLdb.connect("127.0.0.1","eth","eth","eth" )
    cursor = db.cursor()
    sql = "SELECT * FROM Accounts WHERE flag='n' order by Balance desc limit 1"
    cursor.execute(sql)
    results = cursor.fetchall()
    if len(results)!=0:
        acc = results[0][0]
        print "Begin to stat account %s:" % acc
        getTransaction(acc)
        print ""
        
    db.close()
