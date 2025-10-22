package artifacts

/*
Implement DNS artifact generation
- Create A, CNAME, TXT records via Cloudflare API
- Validate records using miekg/dns
- Store records in database
*/