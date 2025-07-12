<a name="Russian"></a>
# Русский
[Jump to English](#English)

<details>
<summary>Содержание</summary>

- [О проекте](#RuAbout)
- [Начало работы](#RuStart)
- [Запуск движка](#RuRun)

</details>

<a name="RuAbout"></a>
### О проекте

Этот проект является приложением к моей ВКР бакалавра. Полный текст моей дипломной работы можно найти на [странице](http://biblio.kosygin-rgu.ru/jirbis2/index.php?option=com_content&view=article&id=41&Itemid=441&vkryear=2025.BAK&page=1&vkrnapr=01.03.02%20%D0%9F%D1%80%D0%B8%D0%BA%D0%BB%D0%B0%D0%B4%D0%BD%D0%B0%D1%8F%20%D0%BC%D0%B0%D1%82%D0%B5%D0%BC%D0%B0%D1%82%D0%B8%D0%BA%D0%B0%20%D0%B8%20%D0%B8%D0%BD%D1%84%D0%BE%D1%80%D0%BC%D0%B0%D1%82%D0%B8%D0%BA%D0%B0) моей университетской электронной библиотеки. Он находится под номером 9 под заголовком "Розин Дмитрий Константинович. Потоко-безопасная система хранения ЗD-объектов с параллельной обработкой их взаимодействия". За данную работу я получил оценку **"Отлично"**.

Это распараллеленный 3D-движок обнаружения и разрешения коллизий, который потенциально может быть использован в качестве основы для симуляционного/игрового 3D-движка. Для упрощения движок в настоящее время поддерживает только сферы в качестве объектов, но количество типов объектов может быть расширено в будущем. Алгоритм обнаружения столкновений основан на комбинации алгоритма **Sweep and Prune** (SaP) для широкой фазы и алгоритма **Separation Axis Theorem** (SAT) для узкой фазы. Алгоритм разрешения столкновений основан на алгоритме **Temporal Gauss-Seidel** (TGS).

<a name="RuStart"></a>
### Начало работы

Подготовка проекта:
- Установите harfang-go
- Исполните `go mod tidy` через командную строку/терминал

Установка harfang-go в Windows:
- Установите msys2
- Исполните `pacman -S mingw-w64-x86_64-gcc` через msys2
- Исполните `setx PATH "%PATH%;C:\msys64\mingw64\bin"` через командную строку/терминал
- Исполните `go get github.com/harfang3d/harfang-go/v3` через командную строку/терминал

<a name="RuRun"></a>
### Запуск движка

Для запуска движка можно либо собрать проект и запустить `.exe`-файл, либо использовать команду `go run main.go` в командной строке/терминале из корневого каталога проекта. В появившемся диалоговом окне выберите варианты алгоритмов, которые вы хотите протестировать. После выбора алгоритмов откроется окно с визуализацией теста. После тестирования между `1024` и `32768` объектами программа автоматически закроется. Из-за ограничений HARFANG окно может зависнуть на большом количестве объектов. В этом случае просто введите `exit` в командной строке/терминале, чтобы закрыть окно, или нажмите `Alt+F4` в экстренных случаях.



<a name="English"></a>
# English
[Перейти к русскому](#Russian)

<details>
  <summary>Contents</summary>
  
  - [About](#EnAbout)
  - [Getting Started](#EnStart)
  - [Running the engine](#EnRun)
  
</details>

<a name="EnAbout"></a>
### About

This project is the codebase appendix for my Bachelor's Thesis. You can find the full text of my Bachelor's Thesis here: [web page](http://biblio.kosygin-rgu.ru/jirbis2/index.php?option=com_content&view=article&id=41&Itemid=441&vkryear=2025.BAK&page=1&vkrnapr=01.03.02%20%D0%9F%D1%80%D0%B8%D0%BA%D0%BB%D0%B0%D0%B4%D0%BD%D0%B0%D1%8F%20%D0%BC%D0%B0%D1%82%D0%B5%D0%BC%D0%B0%D1%82%D0%B8%D0%BA%D0%B0%20%D0%B8%20%D0%B8%D0%BD%D1%84%D0%BE%D1%80%D0%BC%D0%B0%D1%82%D0%B8%D0%BA%D0%B0) of my university e-library. It's listed as entry #9 under the title "Розин Дмитрий Константинович. Потоко-безопасная система хранения ЗD-объектов с параллельной обработкой их взаимодействия". For this work I received an **"Excellent"** grade.

This is a parallelized 3D collision engine, it can potentially be used as a base for a 3D simulation/game engine. For the sake of simplification, the engine currently only supports spheres as objects, the number of object types can be expanded in the future. The collision detection algorithm is based on the combination of the **Sweep and Prune** (SaP) algorithm for the broad phase and **Separation Axis Theorem** (SAT) algorithm for the narrow phase. The collision resolution algorithm is based on the **Temporal Gauss-Seidel** (TGS) algorithm.

<a name="EnStart"></a>
### Getting Started

To prepare the project:
- Install harfang-go
- Run the `go mod tidy` via cmd/terminal

To install harfang-go on windows:
- Install msys2
- Run the `pacman -S mingw-w64-x86_64-gcc` via msys2
- Run the `setx PATH "%PATH%;C:\msys64\mingw64\bin"` via cmd/terminal
- Run the `go get github.com/harfang3d/harfang-go/v3` via cmd/terminal

<a name="EnRun"></a>
### Running the engine

To run the engine you can either build the project and run the `.exe` file or use the `go run main.go` command in the cmd/terminal from the root directory of the project. In the dialogue box that appears, select the variations of the algorithms you want to test. Once you have selected the algorithms, the window will appear with the test visualization. After the test between `1,024` and `32,768` objects the program will close automatically. Due to HARFANG limitations the window may not respond with a large number of objects, in that cases just type `exit` in the cmd/terminal to close the window or use `Alt+F4` in emergency cases.