"""Erzeugt assets/topic-cloud.png, die Wortwolke auf Folie 08b.

Aufruf aus services/slides/ (Python 3 mit: pip install wordcloud pillow fonttools brotli numpy):
    python tools/topic-cloud.py 1
Die Zahl ist der Zufalls-Seed; 1 ergibt die eingebaute Variante. Nach dem
Erzeugen das gewünschte Bild nach assets/topic-cloud.png kopieren.
"""
import sys, random, os
from fontTools.ttLib import TTFont
os.makedirs('out', exist_ok=True)
for w in ('600','800'):
    f=TTFont(f'assets/fonts/poppins-{w}.woff2'); f.flavor=None; f.save(f'out/poppins-{w}.ttf')
from wordcloud import WordCloud
from PIL import Image, ImageDraw, ImageFont
words={
 'Service-Schnitt':8,'Team-Schnitte':7.5,'Datenhoheit':6.5,'Synchron vs. Asynchron':6,'CI/CD & Deployment':6,
 'Auth über Servicegrenzen':5,'Schnittstellen-Versionierung':5,
 'Logging':3.5,'Monitoring':3.5,'Alerting':3,'Contract Testing':3.5,'API-Gateway & BFF':3.5,
 'Security':3,'Messaging':3.5,'Service Mesh':3.5,'Kosten & Betrieb':3.5,
}
KEYS={'Service-Schnitt','Team-Schnitte','Datenhoheit','Synchron vs. Asynchron','CI/CD & Deployment'}
COLORS=['#7348E1','#5B4E93','#0A0349']
W,H=420,420; SC=3
for seed in [int(a) for a in sys.argv[1:]]:
    wc=WordCloud(width=W,height=H,font_path='out/poppins-800.ttf',prefer_horizontal=0.42,margin=10,
                 relative_scaling=0.35,min_font_size=12,max_font_size=42,random_state=seed,
                 background_color=None,mode='RGBA',collocations=False,repeat=False).generate_from_frequencies(words)
    if len(wc.layout_)<len(words):
        print('seed',seed,'dropped',len(words)-len(wc.layout_)); continue
    PAD=10
    img=Image.new('RGBA',((W+2*PAD)*SC,(H+2*PAD)*SC),(246,242,254,255))
    d=ImageDraw.Draw(img)
    rnd=random.Random(seed)
    for (word,freq),fs,(y,x),orient,_ in wc.layout_:
        key=word in KEYS
        font=ImageFont.truetype('out/poppins-800.ttf' if key else 'out/poppins-600.ttf', fs*SC)
        bb=font.getbbox(word)
        tw,th=bb[2]-bb[0],bb[3]-bb[1]
        tile=Image.new('RGBA',(tw+8*SC,th+5*SC),(0,0,0,0))
        td=ImageDraw.Draw(tile)
        if key:
            td.rounded_rectangle([0,0,tile.width-1,tile.height-1],radius=3*SC,fill='#D3F871')
            col='#0A0349'
        else:
            col=rnd.choice(COLORS)
        td.text(((tile.width-tw)//2-bb[0],(tile.height-th)//2-bb[1]),word,font=font,fill=col)
        if orient is not None: tile=tile.rotate(90,expand=True)
        # place centered in the layout box
        lb=ImageFont.truetype('out/poppins-800.ttf',fs).getbbox(word)
        bw,bh=(lb[2]-lb[0],lb[3]-lb[1])
        if orient is not None: bw,bh=bh,bw
        cx=(x+PAD+bw/2)*SC; cy=(y+PAD+bh/2)*SC
        img.alpha_composite(tile,(int(cx-tile.width/2),int(cy-tile.height/2)))
    nv=sum(1 for l in wc.layout_ if l[3] is not None); nk=sum(1 for l in wc.layout_ if l[3] is not None and l[0][0] in KEYS)
    img.save(f'out/cloud-{seed}.png')
    print('seed',seed,'ok vertical',nv,'keys vertical',nk)
