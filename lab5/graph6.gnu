set terminal png size 1000,600
set output "graph6.png"

set title "Переходы цепи Маркова"
set xlabel "Шаг"
set ylabel "Состояние"

set grid
set key right top

plot "graph6_1.dat" using 1:2 with linespoints title "Матрица A", \
     "graph6_2.dat" using 1:2 with linespoints title "Матрица B"