
---

## **Novelty :**

OSINT tools mostly **consume** data (from GitHub, DNS, social feeds, etc.) but we are **feeding** them. It creates harmless DNS records, GitHub commits, repo names, etc., then **analyzes how OSINT systems pick them up, index them, and misinterpret them**.

---

## **Target :**

| Stakeholder                            | Why They Care                                                                                                                  |
| -------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| **Cyber Threat Intel (CTI) Analysts**  | To understand how false indicators distort OSINT feeds and correlation graphs.|
| **OSINT Researchers**                  | To quantify how “truth” propagates - how quickly GitHub commits or DNS records enter public datasets.                          |
| **Privacy & Info Researchers** | To explore disinformation resilience - how adversaries could use similar tactics maliciously, and how to defend against that.  |

---

## **Research POV :**

There’s a ton of work on:

* Adversarial ML (evasion, poisoning)
* Misinformation on social media
* Honeypot and deception systems

But **almost nothing systematic on adversarial OSINT ingestion**.

Cloak provides the infrastructure to study:

* **Signal reliability:** How many OSINT feeds validate what they ingest?
* **Temporal drift:** How long do false signals persist?
* **Cross-source contamination:** Does a false GitHub commit end up referenced in VirusTotal or OTX?

### **Experimental Reproducibility**

Using:

* SQLite
* Audit logs for every API call
* Raw response archival
  means results can be verified and shared without ethical leaks. It’s a **lab-grade testbed**.

### **Quantitative Outputs**

* Visibility rates (% of artifacts detected)
* Detection latency (time to first appearance)
* Persistence curves (time until deletion)
  That turns a qualitative cybersecurity question into a **quantitative, data-science-friendly** one.

### **Cross-domain Insight**

Because it covers **GitHub** (code OSINT) and **DNS** (infra OSINT), it can explore how these two domains **cross-pollinate in intelligence systems**.


# Analogy ->

**Downstream systems (including LLMs that call `/search`) will treat whatever the web returns as “fact” if the ingestion layer doesn’t know how to distrust sources.**

If garbage is present on the internet, automated systems will eat it and regurgitate it unless you (or they) actively check provenance and consistency.

# 1. How the failure happens

1. **Inject false-but-plausible artifact**. e.g a GitHub commit claiming "malware X uses domain y.smthg", or a DNS TXT record that looks like an IOCs list.
2. **OSINT collectors scrape it**.
3. **Aggregators index & correlate**. They mark it as observed.
4. **LLMs workflows query those APIs** or search the web, they see the indexed data and treat it as corroboration.
5. **False positive becomes authoritative** through repetition and cross-references.

LLMs that call `/search` face two problems:

* **Garbage-in**: the search result is wrong (or deliberately planted).
* **Garbage-amplification**: the model presents it with confidence, giving the false signal legitimacy.

# 2. Concrete outcomes to be measured

* **Visibility Rate** = (# of OSINT sources that picked the artifact) / (total monitored sources).
* **Propagation Latency** = time(first index in any monitored feed) − time(injection).
* **Contamination Depth** = number of distinct systems that reference the artifact (e.g., GitHub → VirusTotal → MISP event → blog).
* **Persistence** = time until artifact is removed / corrected (or infinity if never removed).

ex formulas:

* VisibilityRate = detected_sources / monitored_sources
* MeanLatency = average(time_detected_i − time_injected)
* PersistenceHalfLife = time by which 50% of sources have removed or corrected it

# 3. How to test it

1. **Design safe artifacts** : dummy repo names (creating a new github account works), commits, DNS TXT records containing clear markers (e.g., `CLOAK_TEST_ABC123`) so they’re non-actionable.
2. **Inject in controlled channels** : create GitHub commits, post a benign paste, publish an innocuous domain/record on a sandbox DNS provider.
3. **Log everything** : local audit (timestamps, content, unique hashes).
4. **Monitor a fixed set of public feeds** at intervals (GitHub, VirusTotal, OTX, Shodan).
5. **Archive all hits** : (save raw responses).
6. **Map propagation graph** : record where it appeared, when, and which downstream references cite it.
7. **Analyze metrics**.

# 4. DB schema (SQLite)

```sql
CREATE TABLE artifacts (
  id TEXT PRIMARY KEY,
  type TEXT,
  content TEXT,
  injected_at DATETIME
);

CREATE TABLE observations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  artifact_id TEXT,
  source TEXT,
  observed_at DATETIME,
  raw_response TEXT,
  FOREIGN KEY(artifact_id) REFERENCES artifacts(id)
);
```

# 5. Defenses & best practises

* **Provenance scoring:** track origin, author history, domain age, hosting provider, and assign trust score.
* **Cross-source sanity checks:** don’t treat a single source as truth, `require > 1` independent corroboration with different provenance.
* **Human supervision:** automated correlation should be suggestions, not final.

---
