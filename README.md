# Midi Manipulator

## Domain-specific declarative language specification for configuration of MIDI-device controls

### Иерархия
```
midi_devices:
  - device_name: MPD226
    startup_delay: 100
    reconnect_interval: 2000
    active: true
    hold_delta: 1000
    namespace: default
    accumulate_controls:
      - keys:
          - 16
          - 17
        rotate: false
        value_range:
          - 0
          - 127
        initial_value: 0
        triggers: null
        increment: 127
        decrement: 1
      - keys:
          - 18
          - 19
        rotate: false
        value_range:
          - 63
          - 127
        initial_value: 0
        triggers: null
        increment: 1
        decrement: 127

  ...
```
### Атрибуты

#### midi_devices 

Тип аргументов: Array   
   
Описание: В данной секции необходимо перечислить все используемые устройства. Их количество должно совпадать с реальным набором устройств подключенных к хабу пультов, иначе сервис прекратит свою работу с ошибкой.

#### device_name 

Тип аргументов: String   
   
Описание: В данной секции необходимо указать название MIDI порта как он определен в системе. **Если указанного порта не существует в системе выполнение сервиса закончится ошибкой.**

Список доступных портов MIDI устройств в системе можно узнать так
```
amidi -l
```
или через список всех USB устройств в системе
```
lsusb -v 2>&1 | grep -e iProduct
```

#### startup_delay 

Тип аргументов: Integer   
   
Описание: Параметр позволяющий указать **строго положительное** время задержки **(в миллисекундах)** до воспроизведения стартовой подсветки устройства. **Рекомендуется указывать если при подключении устройства наблюдается конфликт заводской подсветки с текущей пользовательской конфигурацией подсветки.**

#### reconnect_interval 

Тип аргументов: Integer   
   
Описание: Параметр позволяющий указать **строго положительное** время между попытками переподключения к недоступному на текущий момент устройству из конфигурации **(в миллисекундах)**.

Ограничения: > 1000 ms.

#### active 

Тип аргументов: Boolean   
   
Описание: Указание данного параметра позволяет определить активно или нет данное MIDI устройство. В активном состоянии процесс считывания сигналов и выполнения команд производится в обычном режиме. В неактивном состоянии приостанавливается процесс считывания сигналов с устройства.

#### hold_delta 

Тип аргументов: Integer   
   
Описание: Служит для определения времени удержания клавиши (в мс) для отправки соответствующего сигнала.

#### namespace 

Тип аргументов: String   
   
Описание: Служит для определения текущей раскладки по названию.

#### accumulate_controls 

Тип аргументов: Struct[]   
   
Описание: Список множеств элементов управления устройства, которые вовзращают неполный диапазон значений velocity (стандартный диапазон энкодеров/потенциометров [0, 127]). Каждое из множеств может использоваться для обобщения логики работы нескольких отдельных элементов управления.

#### keys 

Тип аргументов: IntArray   
   
Описание: Набор id элементов управления устройства для которых требуется определить особые правила обработки сигналов.

Ограничения: **Пересечение id элементов управления с другими множествами одного устройства приведет к переопределению их логики.**

#### rotate 

Тип аргументов: Boolean   
   
Описание: Способ обработки сигналов в зависимости от возвращаемых результатов после вращения.

#### value_range

Тип аргументов: IntArray[2]   
   
Описание: Переопределение границ диапазона доступных значений для указанных элементов управления.

#### initial_value

Тип аргументов: Integer   
   
Описание: Начальное значение velocity независимо от текущего положения элемента управления.

#### triggers

Тип аргументов: Struct   
   
Описание: Набор значений при достижении которых происходит изменение текущего значения velocity.

#### increment/decrement

Тип аргументов: Integer   
   
Описание: Значение полученное с устройства при котором происходит инкрементация/декрементация текущего значения velocity.

## Domain-specific declarative language specification for backlight configuration of MIDI devices

### Иерархия
```
device_light_configuration:
  - device_name: MPD226
    backlight_time_offset: 50
    color_spaces:
      - color_space_id: 1
        on:
          - color_name: red
            payload: 7F
        off:
          - color_name: black
            payload: '00'
    keyboard_backlight:
      - key_range:
          - 0
          - 3
        key_number_shift: 0
        color_space: 1
        statuses:
          on:
            type: ControlChange
            fallback_color: red
            bytes: B1 %key %payload
          off:
            type: ControlChange
            fallback_color: black
            bytes: B1 %key %payload
      - key_range:
          - 60
          - 88
        key_number_shift: 0
        color_space: 1
        statuses:
          on:
            type: NoteOn
            fallback_color: red
            bytes: 91 %key %payload
          off:
            type: NoteOn
            fallback_color: black
            bytes: 81 %key %payload
```
### Атрибуты

#### device_light_configuration 

Тип аргументов: Array   
   
Описание: Список конфигураций подсветки для каждого известного MIDI устройства

#### device_name 

Тип аргументов: String   
   
Описание: Название MIDI порта. **Должно совпадать с аналогичным названием в user config для одного устройства.**

#### backlight_time_offset 

Тип аргументов: Integer   
   
Описание: Временной промежуток (в мс) между подсветкой клавиш.

#### color_spaces

Тип аргументов: Array   
   
Описание: Набор цветовых палитр для разных диапазонов клавиш.

#### color_space_id 

Тип аргументов: Integer   
   
Описание: Id отдельного color_space.

#### on/off 

Тип аргументов: Array   
   
Описание: Массив доступных цветов подсветки для набора клавиш в включенном/выключенном состоянии.

#### color_name

Тип аргументов: String   
   
Описание: Название цвета.

#### payload 

Тип аргументов: String   
   
Описание: Массив байтов определяющий цвет для конкретного набора клавиш, который представен в строковом виде.

#### keyboard_backlight

Тип аргументов: Array   
   
Описание: Массив диапазонов клавиш устройства для которых возможно управление подсветкой.

#### key_range

Тип аргументов: IntArray[2]   
   
Описание: Диапазон клавиш для которых подсветка работает по одному принципу. **Диапазоны не должны пересекаться, иначе будет undefined behaviour.**

#### key_number_shift

Тип аргументов: Integer   
   
Описание: Сдвиг id клавиш.

#### color_space

Тип аргументов: Integer   
   
Описание: Привязка к цветовой палитре для текущего диапазона клавиш.

#### statuses

Тип аргументов: Struct   
   
Описание: Принцип работы подсветки клавиш во включенном/выключенном состоянии.

#### on/off 

Тип аргументов: Struct   
   
Описание: Включенное/выключенное состоняние.

#### type

Тип аргументов: String   
   
Описание: Тип MIDI сообщения для подсветки.

#### fallback_color

Тип аргументов: String   
   
Описание: Название базового цвета, если не были указаны другие или их не существует в палитре цветов для данного диапазона. **Данный цвет обязательно должен быть в палитре для указанного диапазона.**

#### bytes

Тип аргументов: String   
   
Описание: Массив байтов в строковом виде для активации подсветки, который содержит ключи форматирования для динамической вставки id клавиши (%key) и набора байтов для указания цвета из палитры (%payload). 
