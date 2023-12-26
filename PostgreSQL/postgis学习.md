---
tags:
  - 数据库
  - gis
date created: 2023-11-13 02:41
date modified: 2023-11-13 02:56
---
[PostGIS官方教程汇总目录_postgis教程-CSDN博客](https://blog.csdn.net/qq_35732147/article/details/85256640)

加载 PostGIS 空间扩展
```sql
create extension postgis;
```

查看是否安装扩展
```sql
select postgis_full_version();
```

# 空间函数

ST_AsText 把 geom 类型转为文本
```sql
select ST_AsText(geom) from streets where name = 'point';
```

## 点空间函数

- ST_X(geometry)：返回 X 坐标
- ST_Y(geometry)：返回 Y 坐标

## 线段空间函数

- ST_Length(geom)：返回线串长度
- ST_StartPoint(geom)：将线串的第一个坐标作为点返回
- ST_EndPoint(geom)：将线串的最后一个坐标作为点返回
- ST_NPoints(geom)：返回线串的点数量

## 多边形空间函数

- ST_Area(geom)：返回多边形面积
- ST_NRings(geom)：返回多边形中环的数量
- ST_ExteriorRing(geom)：以线串的形式返回多边形最外面的环
- ST_InteriorRingN(geom)：以线串形式返回指定的内部环
- ST_Perimeter(geom)：返回所有环的长度

## 集合

- MultiPoint：点集合
- MultiLineString：线串集合
- MultiPolygon：多边形集合
- GeometryCollection：由任意几何图形组成的异构集合

## 集合空间函数

- ST_NumGeometries(geom)：返回集合中的元素数量
- ST_GeometryN(geom,n)：返回几何集合的一个元素
- ST_Area(geom)：返回多边形几何体的面积
- ST_Length(geom)：返回线性几何体的二维长度