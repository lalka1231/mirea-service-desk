import java.util.Arrays;
import java.util.Scanner;

public class Main {
    public static void main(String[] args) {
        System.out.println("=== Проверка работы основных конструкций ===");

        // 1. Проверка объявления переменных и операций
        System.out.println("\n1. Переменные и операции:");
        int a = 10;
        int b = 5;
        System.out.println("a + b = " + (a + b));
        System.out.println("a * b = " + (a * b));

        // 2. Проверка работы с массивом
        System.out.println("\n2. Работа с массивом:");
        int[] numbers = {3, 7, 2, 9, 5};
        System.out.println("Массив: " + Arrays.toString(numbers));
        System.out.println("Длина массива: " + numbers.length);
        System.out.println("3-й элемент массива: " + numbers[2]);

        // 3. Проверка условной конструкции
        System.out.println("\n3. Условная конструкция (if-else):");
        if (a > b) {
            System.out.println("a больше b");
        } else {
            System.out.println("a не больше b");
        }

        // 4. Проверка циклов
        System.out.println("\n4. Циклы:");

        // Цикл for
        System.out.print("Цикл for: ");
        for (int i = 0; i < 3; i++) {
            System.out.print(i + " ");
        }

        // Цикл while
        System.out.print("\nЦикл while: ");
        int count = 0;
        while (count < 3) {
            System.out.print(count + " ");
            count++;
        }

        // Цикл do-while
        System.out.print("\nЦикл do-while: ");
        count = 0;
        do {
            System.out.print(count + " ");
            count++;
        } while (count < 3);

        // 5. Проверка ввода данных
        System.out.println("\n\n5. Ввод данных с клавиатуры:");
        Scanner scanner = new Scanner(System.in);
        System.out.print("Введите ваше имя: ");
        String name = scanner.nextLine();
        System.out.println("Привет, " + name + "!");

        // 6. Проверка работы со строками
        System.out.println("\n6. Методы класса String:");
        String testString = "Hello Java";
        System.out.println("Исходная строка: " + testString);
        System.out.println("Длина строки: " + testString.length());
        System.out.println("В верхнем регистре: " + testString.toUpperCase());
        System.out.println("Замена 'a' на 'o': " + testString.replace('a', 'o'));
        System.out.println("Индекс символа 'J': " + testString.indexOf('J'));

        // 7. Проверка метода
        System.out.println("\n7. Работа метода:");
        int number = 5;
        long factorial = calculateFactorial(number);
        System.out.println("Факториал числа " + number + " = " + factorial);

        System.out.println("\n=== Проверка завершена ===");
    }

    // Метод для вычисления факториала
    public static long calculateFactorial(int n) {
        long result = 1;
        for (int i = 1; i <= n; i++) {
            result *= i;
        }
        return result;
    }
}