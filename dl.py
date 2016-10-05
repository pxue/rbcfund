import os, json, urllib2

FILE = "http://funds.rbcgam.com/pdf/fund-pages/monthly/%s_e.pdf"

D = {}
with open("./data/fundList.json") as f:
    for f in json.loads(f.read())["Data"]:
        fcode = f["FundCode"]
        fname = f["FundName"]
        cat = f["AssetClass"]
        if not D.get(cat):
            D[cat] = []

        if "US$" in fname:
            continue

        #if "Equity" in f["MorningstarCategoryName"] or \
                #"Equity" in fname:
            #continue

        D[cat].append((fcode, fname))

        #fname = "portfolios/%s.pdf" % fcode.upper()
        #if os.path.isfile(fname):
            #continue

        #urlpath = FILE % fcode
        #try:
            #pdf = urllib2.urlopen(urlpath)
            #fsave = open(fname, "wb")
            #fsave.write(pdf.read())
            #pdf.close()
            #fsave.close()

            #print "saved", fcode
        #except Exception as e:
            #print e, f["FundName"]

HTML = "<html><body>%s</body></html>"
BODY = ""

for k, v in D.iteritems():
    BODY += "<h1>%s</h1>" % k
    BODY += "<ul>"
    for fund in v:
        BODY += "<li>"
        BODY += "<a href='./data/portfolios/%s.pdf' target='_blank'>%s - %s</a>" % (fund[0], fund[0], fund[1])
        BODY += "</li>"
    BODY += "</ul>"

HTML = HTML % BODY
print HTML
