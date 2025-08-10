(function(){
  function qs(sel, root=document){ return root.querySelector(sel); }
  function qsa(sel, root=document){ return Array.from(root.querySelectorAll(sel)); }

  function setYear(){ const y=qs('#year'); if (y) y.textContent = new Date().getFullYear(); }

  function smoothAnchors(){
    qsa('a[href^="#"]').forEach(a => {
      a.addEventListener('click', (e) => {
        const id = a.getAttribute('href');
        if (!id || id === '#' || id.length < 2) return;
        const el = qs(id);
        if (el) {
          e.preventDefault();
          el.scrollIntoView({behavior:'smooth', block:'start'});
        }
      });
    });
  }

  function navToggle(){
    const btn = qs('#navToggle');
    const menu = qs('#menu');
    if (!btn || !menu) return;
    btn.addEventListener('click', () => {
      const open = menu.getAttribute('data-open') === 'true';
      menu.setAttribute('data-open', String(!open));
      btn.setAttribute('aria-expanded', String(!open));
    });
    // close menu on link click (mobile)
    qsa('#menu a').forEach(link => link.addEventListener('click', () => {
      if (window.innerWidth <= 820) {
        menu.setAttribute('data-open', 'false');
        btn.setAttribute('aria-expanded', 'false');
      }
    }));
  }

  function revealOnScroll(){
    const items = qsa('.reveal');
    if (!('IntersectionObserver' in window) || items.length === 0) {
      items.forEach(el => el.setAttribute('data-in','true'));
      return;
    }
    const io = new IntersectionObserver((entries)=>{
      entries.forEach(e => {
        if (e.isIntersecting){ e.target.setAttribute('data-in','true'); io.unobserve(e.target); }
      })
    }, {threshold: 0.12});
    items.forEach(el => io.observe(el));
  }

  // Hearts
  const MAX_HEARTS = 60;
  let heartCount = 0;
  function spawnHeart(x, y){
    if (heartCount > MAX_HEARTS) return;
    const h = document.createElement('div');
    h.className = 'heart';
    h.style.left = (x - 7) + 'px';
    h.style.top = (y - 7) + 'px';
    h.style.transform += ` translateX(${(Math.random()*40-20)}px)`;
    document.body.appendChild(h);
    heartCount++;
    setTimeout(()=>{ h.remove(); heartCount--; }, 2600);
  }
  function burstHeartsAt(el){
    const rect = el.getBoundingClientRect();
    const cx = rect.left + rect.width/2 + window.scrollX;
    const cy = rect.top + rect.height/2 + window.scrollY;
    for (let i=0;i<12;i++) setTimeout(()=>spawnHeart(cx + (Math.random()*60-30), cy + (Math.random()*20-10)), i*60);
  }

  function enableAmbientHearts(){
    document.addEventListener('click', (e)=>{
      spawnHeart(e.pageX, e.pageY);
    });
  }

  function handleRSVP(){
    const form = qs('#rsvpForm');
    const notice = qs('#rsvpNotice');
    if (!form) return;
    form.addEventListener('submit', async (e)=>{
      e.preventDefault();
      const fd = new FormData(form);
      const name = (fd.get('name')||'').toString().trim();
      const email = (fd.get('email')||'').toString().trim();
      const attending = (fd.get('attending')||'').toString().trim();
      if (!name || !email || !attending) {
        alert('Please fill in name, email and attendance.');
        return;
      }
      // Here we could POST to an API endpoint; for now simulate success
      await new Promise(r=>setTimeout(r, 500));
      if (notice) notice.setAttribute('data-show','true');
      burstHeartsAt(form);
      form.reset();
    });
  }

  function inviteCard(){
    const card = qs('#inviteCard');
    if (!card) return;
    card.addEventListener('click', ()=>{
      const open = card.getAttribute('aria-expanded') === 'true';
      card.setAttribute('aria-expanded', String(!open));
    });
  }

  // Invitation page handling (code entry + valid invitation view)
  function setupInvitationPage(){
    const ctx = qs('#inviteCtx');
    const form = qs('#inviteCodeForm');
    const input = qs('#inviteCodeInput');
    const closeBtn = qs('#closeInvite');
    const counterEl = qs('#acceptedCounter');
    const myTable = qs('#myInvitationTable');
    const allTable = qs('#allAcceptedTable');

    // Early out if not on invitation page at all
    if (!ctx && !form) return;

    const code = ctx ? (ctx.getAttribute('data-code') || '') : '';
    const valid = ctx ? (ctx.getAttribute('data-valid') === 'true') : false;

    // When valid invitation page, store current invite code
    if (valid && code) {
      try { localStorage.setItem('invite:current', code); } catch(e){}
    }

    if (closeBtn){
      closeBtn.addEventListener('click', ()=>{
        try {
          const current = localStorage.getItem('invite:current');
          // clear stored per-invite statuses
          Object.keys(localStorage).forEach(k=>{ if (k.startsWith('inviteStatus:')) localStorage.removeItem(k); });
          localStorage.removeItem('invite:current');
        } catch(e){}
        window.location.href = '/bryllup/invitation/';
      });
    }

    if (form) {
      form.addEventListener('submit', (e)=>{
        e.preventDefault();
        const c = (input && input.value || '').trim();
        if (!c) { alert('Indtast venligst din invitationskode.'); return; }
        window.location.href = '/bryllup/invitation/' + encodeURIComponent(c) + '/';
      });
    }

    const weddingDate = new Date('2025-09-20T13:00:00'); // local time
    const lockFrom = new Date(weddingDate.getTime() - 7*24*60*60*1000);

    function locked(){ return new Date() >= lockFrom; }

    function renderMyTable(data){
      if (!myTable) return;
      const tbody = myTable.querySelector('tbody');
      if (!tbody) return;
      tbody.innerHTML = '';
      const acceptedSet = new Set(data.accepted||[]);
      (data.members||[]).forEach(name => {
        const tr = document.createElement('tr');
        const tdName = document.createElement('td'); tdName.textContent = name;
        const tdStatus = document.createElement('td');
        const tdAct = document.createElement('td');
        const accBtn = document.createElement('button'); accBtn.className='button'; accBtn.textContent='Deltager';
        const decBtn = document.createElement('button'); decBtn.className='button'; decBtn.style.marginLeft='8px'; decBtn.textContent='Deltager ikke';
        const isLocked = locked();
        const declinedKey = `inviteDeclined:${code}:${name}`;
        const isAccepted = acceptedSet.has(name);
        const isDeclined = !isAccepted && localStorage.getItem(declinedKey) === '1';
        // status pill
        const statusSpan = document.createElement('span');
        statusSpan.className = 'status';
        if (isAccepted) { statusSpan.classList.add('accepted'); statusSpan.textContent = 'Deltager'; }
        else if (isDeclined) { statusSpan.classList.add('declined'); statusSpan.textContent = 'Deltager ikke'; }
        else { statusSpan.classList.add('pending'); statusSpan.textContent = 'Afventer'; }
        tdStatus.appendChild(statusSpan);
        // button styles
        if (isAccepted) { accBtn.classList.add('primary'); decBtn.classList.remove('danger'); }
        else if (isDeclined) { decBtn.classList.add('danger'); accBtn.classList.remove('primary'); }
        // disable per rules
        accBtn.disabled = !valid || isLocked;
        decBtn.disabled = !valid;
        accBtn.title = isLocked ? 'Accepter er låst 7 dage før' : '';
        accBtn.addEventListener('click', async ()=>{
          try { localStorage.removeItem(declinedKey); } catch(e){}
          await postJSON(`/bryllup/api/invites/${encodeURIComponent(code)}/accept/`, {name});
          await refreshAllAccepted();
        });
        decBtn.addEventListener('click', async ()=>{
          try { localStorage.setItem(declinedKey, '1'); } catch(e){}
          await postJSON(`/bryllup/api/invites/${encodeURIComponent(code)}/decline/`, {name});
          await refreshAllAccepted();
        });
        tdAct.appendChild(accBtn); tdAct.appendChild(decBtn);
        tr.appendChild(tdName); tr.appendChild(tdStatus); tr.appendChild(tdAct);
        tbody.appendChild(tr);
      });
    }

    function renderAllTable(list){
      if (!allTable) return;
      const tbody = allTable.querySelector('tbody');
      if (!tbody) return;
      tbody.innerHTML = '';
      const arr = Array.isArray(list) ? list : [];
      arr.forEach(name => {
        const tr = document.createElement('tr');
        const tdName = document.createElement('td'); tdName.textContent = name;
        tr.appendChild(tdName);
        tbody.appendChild(tr);
      });
      const title = document.getElementById('allAcceptedTitle');
      if (title){
        const n = arr.length || tbody.querySelectorAll('tr').length;
        title.textContent = `Deltagerlisten (${n} har klikket deltager)`;
      }
    }

    async function refreshAccepted(){
      try {
        const res = await fetch(`/bryllup/api/invites/${encodeURIComponent(code)}/accepted/`);
        const data = await res.json();
        if (counterEl) counterEl.textContent = data.count + '/' + data.capacity;
        renderMyTable(data);
      } catch(e) { /* ignore for now */ }
    }

    async function refreshAllAccepted(){
      try {
        const res = await fetch('/bryllup/api/invites/accepted/');
        const data = await res.json();
        renderAllTable(data.accepted||[]);
      } catch(e) { /* ignore */ }
    }

    async function postJSON(url, payload){
      await fetch(url, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(payload)});
      await refreshAccepted();
    }

    if (valid){
      refreshAccepted();
      refreshAllAccepted();
    }
  }

  // Calendar helpers
  function pad(n){ return String(n).padStart(2,'0'); }
  function toUTCStringBasic(d){
    return d.getUTCFullYear()+''+pad(d.getUTCMonth()+1)+''+pad(d.getUTCDate())+'T'+pad(d.getUTCHours())+pad(d.getUTCMinutes())+pad(d.getUTCSeconds())+'Z';
  }
  function setupCalendarPage(){
    const a = qs('#googleCal');
    const btn = qs('#icsBtn');
    if (!a && !btn) return;
    const start = new Date('2025-09-20T13:00:00');
    const end = new Date('2025-09-20T23:59:00');
    const text = 'Bryllup – Lars & Ida';
    const details = 'Vi fejrer kærligheden! Se praktisk info på /bryllup/info/';
    const location = 'Farre Kirke & Fælleshuset';

    if (a){
      const href = 'https://www.google.com/calendar/render?action=TEMPLATE'
        + '&text=' + encodeURIComponent(text)
        + '&dates=' + toUTCStringBasic(start) + '/' + toUTCStringBasic(end)
        + '&details=' + encodeURIComponent(details)
        + '&location=' + encodeURIComponent(location)
        + '&sf=true&output=xml';
      a.href = href;
    }

    if (btn){
      btn.addEventListener('click', ()=>{
        const uid = 'wedding-'+Date.now()+'@localhost';
        const ics = [
          'BEGIN:VCALENDAR',
          'VERSION:2.0',
          'PRODID:-//Vores Bryllup//DA',
          'CALSCALE:GREGORIAN',
          'METHOD:PUBLISH',
          'BEGIN:VEVENT',
          'UID:'+uid,
          'DTSTAMP:'+toUTCStringBasic(new Date()),
          'DTSTART:'+toUTCStringBasic(start),
          'DTEND:'+toUTCStringBasic(end),
          'SUMMARY:'+text,
          'DESCRIPTION:'+details.replace(/\n/g,'\\n'),
          'LOCATION:'+location,
          'END:VEVENT',
          'END:VCALENDAR' ].join('\r\n');
        const blob = new Blob([ics], {type:'text/calendar'});
        const url = URL.createObjectURL(blob);
        const tmp = document.createElement('a');
        tmp.href = url; tmp.download = 'wedding.ics';
        document.body.appendChild(tmp); tmp.click(); tmp.remove();
        setTimeout(()=>URL.revokeObjectURL(url), 2000);
      });
    }
  }

  function setupConsentPopup(){
    try {
      if (localStorage.getItem('consent:activity') === '1') return; // already accepted
    } catch(e) { /* ignore storage errors */ }

    const box = document.createElement('div');
    box.className = 'cookie-popup';
    box.setAttribute('role', 'dialog');
    box.setAttribute('aria-live', 'polite');
    box.innerHTML = [
      '<p>Vi bruger ikke cookies, men aktivitet på siden bliver logget. Vil du acceptere det?</p>',
      '<div class="actions">',
      '  <button type="button" class="button primary" id="consentAccept">Accepter</button>',
      '  <button type="button" class="button secondary" id="consentReject">Afvis</button>',
      '</div>',
      '<div class="small" id="consentMsg" style="display:none;margin-top:6px;"></div>'
    ].join('');

    document.body.appendChild(box);
    requestAnimationFrame(()=>{ box.setAttribute('data-show','true'); });

    const acceptBtn = box.querySelector('#consentAccept');
    const rejectBtn = box.querySelector('#consentReject');
    const msg = box.querySelector('#consentMsg');

    function hide(){ box.removeAttribute('data-show'); setTimeout(()=>box.remove(), 250); }

    if (acceptBtn){
      acceptBtn.addEventListener('click', ()=>{
        try { localStorage.setItem('consent:activity','1'); } catch(e) {}
        hide();
      });
    }
    if (rejectBtn){
      rejectBtn.addEventListener('click', ()=>{
        if (msg){
          msg.textContent = 'Jamen så lad da være med at bruge siden!';
          msg.style.display = 'block';
        }
        // Disable buttons after reject message
        if (acceptBtn) acceptBtn.disabled = true;
        rejectBtn.disabled = true;
        setTimeout(()=>{ hide(); }, 5000);
      });
    }
  }

  document.addEventListener('DOMContentLoaded', function(){
    setYear();
    smoothAnchors();
    navToggle();
    revealOnScroll();
    enableAmbientHearts();
    handleRSVP();
    inviteCard();
    setupInvitationPage();
    setupCalendarPage();
    setupConsentPopup();
  });
})();