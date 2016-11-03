Action
======

An action to perform on either the selected element or the page state.


name
----

================= ==========
**required**      false
**type**          string
**default**       ""
================= ==========

The name of the action. This will be output for reference of the current action being run.

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 1
   :lines: 22-24


type
----

The type of action to run.

Options
^^^^^^^

================= ==
Type
================= ==
click             ✓
fill              ✓
check             ✓
clear             ✓
double_click      ✓
flick_finger      ✓
mouse_to_element  ✓
scroll_finger     ✓
select            ✓
submit            ✓
tap               ✓
touch             ✓
uncheck           ✓
upload_file       ✓
send_keys         ✓
================= ==

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 2
   :lines: 22-24


text
----

================= ==========
**required**      false
**type**          string
**default**       ""
================= ==========

The text to send to the event.

Supported Types
^^^^^^^^^^^^^^^

========== ==
Type
========== ==
fill       ✓
send_keys  ✓
========== ==

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 3
   :lines: 22-24


assertions
----------

================= ====================
**required**      false
**type**          [] :doc:`assertion`
**default**       nil
================= ====================

The assertions to run against alongside the action.
