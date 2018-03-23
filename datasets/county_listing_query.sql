-- After importing shapefile with:
-- shp2pgsql -W LATIN1 -s 4269 ./datasets/gz_2010_us_050_00_500k.shp public.county_shapes | psql -h localhost -d counties

SELECT
       t.gid,
       geo_id,
       state,
       county,
       name,
       lsad,
       round(max(ST_Distance(dump.geom::geography, ST_Centroid(t.geom)::geography))/1000) AS max_radius,
       ST_X(ST_Centroid(t.geom)) as lon,
       ST_y(ST_Centroid(t.geom)) as lat
FROM county_shapes t
JOIN ST_DumpPoints(t.geom) dump ON true
where state in ('01','04','05','06','08','09','10','12','13','16','17','18','19','20','21','22','23','24','25','26','27','28','29','30','31','32','33','34','35','36','37','38','39','40','41','42','44','45','46','47','48','49','50','51','53','54','55','56')
        and lsad in ('County', 'Parish')
GROUP BY 1, 2, 3, 4, 5, 6
ORDER BY 1
;
