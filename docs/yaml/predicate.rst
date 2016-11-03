Predicate
=========

A predicate is used to select one or more element on the page.


name
----

================= ==========
**required**      false
**type**          string
**default**       ""
================= ==========

The name of the predicate. This will be output for reference of the current predicate being run.

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 2
   :lines: 7-11


type
----

================= ==========
**required**      false
**type**          string
**default**       "default"
================= ==========

The type of predicate.

Options
^^^^^^^

=========== ============== ========= =========
Type        find (default) multi     first
=========== ============== ========= =========
default     ✓              ✓         ✓
by_button   ✓              ✓         ✓
by_class    ✓              ✓         ✓
by_id       ✓              ✓
by_label    ✓              ✓         ✓
by_link     ✓              ✓         ✓
by_name     ✓              ✓         ✓
by_xpath    ✓              ✓         ✓
for_appium  ✓
=========== ============== ========= =========

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 3
   :lines: 7-11


selector
--------

================= ==========
**required**      true
**type**          string
**default**       ""
================= ==========

The selector for finding the element.

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 4
   :lines: 7-11


text
----

================= ======================
**required**      type == ``by_appium``
**type**          string
**default**       ""
================= ======================

The text to select with. Only used for ``by_appium`` type.

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 4
   :lines: 12-15


multi
-----

================= ==========
**required**      false
**type**          boolean
**default**       false
================= ==========

Whether or not the predicate should do multi-selection.

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 5
   :lines: 7-11


first
-----

================= ==========
**required**      false
**type**          boolean
**default**       false
================= ==========

Whether or not to only return the first element.

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 4
   :lines: 16-19
