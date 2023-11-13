update ap_trust_submissions apts2  set submitted_at =
(select submitted_at from ap_trust_statuses apts1 where apts2.metadata_id = apts1.metadata_id);