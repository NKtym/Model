set terminal pngcairo size 1600,700 enhanced font 'Arial,14'
set output 'hist_density2.png'
set title 'Гистограмма выборки (нормированная) и теоретическая плотность f(x)'
set xlabel 'x'
set ylabel 'Плотность'
set grid
set xrange [-0.2:2.5]
set yrange [0:*]
set xtics -0.2,0.1,2.5
set tics scale 0.6
set xtics out
set grid xtics    
set style fill solid 0.7 border -1
plot 'hist.dat' using 1:2 with boxes lc rgb 'gray' title 'Нормированная гистограмма', \
     'density.dat' using 1:2 with lines lt 1 lw 3 lc rgb 'red' title 'f(x)'
