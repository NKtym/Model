set terminal png size 1000,600
set output "graph7.png"

set title "Автокорреляция цепи Маркова"
set xlabel "Лаг"
set ylabel "Автокорреляция"

set grid
set key right top

plot "graph7_1.dat" using 1:2 with linespoints title "Матрица A", \
     "graph7_2.dat" using 1:2 with linespoints title "Матрица B"