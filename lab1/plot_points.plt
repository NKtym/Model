set terminal pngcairo size 1600,700 enhanced font 'Arial,14'
set output 'density_points.png'
set title 'Плотность f(x) и точки выборки (каждая точка на координате (x, f(x)))'
set xlabel 'x'
set ylabel 'Плотность'
set grid
set xrange [-0.3:2.6]
set xtics -0.3,0.1,2.6
set tics scale 0.6
set xtics out
set yrange [-0.3:1.3]
set ytics -0.3,0.2,1.3
set xzeroaxis lt 1 lw 2 lc rgb 'grey'
set yzeroaxis lt 1 lw 2 lc rgb 'grey'
plot 'density.dat' using 1:2 with lines lt 1 lw 3 lc rgb 'blue' title 'f(x)', \
     'samples_on_density.dat' using 1:2 with points pt 7 ps 0.6 lc rgb 'black' title 'Точки выборки'
