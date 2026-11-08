<!-- .slide: data-background-color="#0A0349" -->

<div class="intro">

<p class="intro-kicker">Microservices richtig gemacht</p>
<h1 class="intro-title">Microservices</h1>
<p class="intro-claim">&hellip; sind auch nur Monolithen mit <span class="intro-hl">Netzwerkproblemen</span></p>

<div class="intro-bottom">
<div class="speaker">
<div>
<div class="speaker-name">Lars Michaelis</div>
<div class="speaker-role">Software-Architekt</div>
</div>
<ul class="speaker-links">
<li><a href="mailto:l.michaelis@neusta.de">l.michaelis@neusta.de</a></li>
<li><a href="https://github.com/larmic/workshop_microservices" target="_blank" rel="noopener">github.com/larmic/workshop_microservices</a></li>
<li><a href="https://www.linkedin.com/in/i-am-lars-michaelis" target="_blank" rel="noopener">linkedin.com/in/i-am-lars-michaelis</a></li>
<li><a href="https://www.youtube.com/@LarsMichaelis" target="_blank" rel="noopener">youtube.com/@LarsMichaelis</a></li>
</ul>
</div>
<div class="intro-trace" role="img" aria-label="Trace einer Buchung: car lehnt ab, flight wird kompensiert">
<div class="trace-meta">trace_id 4bf92f3577b34da6 &middot; POST /booking/bookings</div>
<div class="trace-rows">
<span>booking</span><div class="trace-bar booking"></div>
<span>flight</span><div class="trace-bar flight"></div>
<span>hotel</span><div class="trace-bar hotel"></div>
<span>car</span><div class="trace-fail"><div class="trace-bar"></div><span>409</span></div>
<span>flight</span><div class="trace-bar compensate"></div>
</div>
<div class="trace-meta">car lehnt ab &rarr; flight wird kompensiert</div>
</div>
</div>

</div>

Note:
- <strong>Begr&uuml;&szlig;ung:</strong> Willkommen bei &bdquo;Microservices richtig gemacht&ldquo;. Kurz zu mir: Rolle, und warum mich das Thema umtreibt.
- <strong>Untertitel:</strong> Die Provokation ist Absicht. Nur stehen lassen, wir kl&auml;ren sie sp&auml;ter in der Diskussion.
- <strong>Trace rechts:</strong> Nicht erkl&auml;ren, nur ank&uuml;ndigen: &bdquo;So sieht eine Buchung am Ende von Tag 2 aus. Bis dahin wisst ihr, was jeder Balken bedeutet.&ldquo;
- <strong>Vorstellungsrunde:</strong> Name, Rolle, eine Erwartung. Erwartungen aufs Flipchart, am Ende gleichen wir ab.
- <strong>&Uuml;berleitung:</strong> Bevor es losgeht, kurz zum Feedback.
