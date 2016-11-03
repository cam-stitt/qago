Case
====

A suite is a singular test case to be run.


name
----

================= ==========
**required**      false
**type**          string
**default**       ""
================= ==========

The name of the case. This will be output for reference of the current test being run.

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 1
   :lines: 1-3


location
--------

================= ==========
**required**      true
**type**          string
**default**       ""
================= ==========

The url location to initially load for the page.

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 2
   :lines: 1-3


browser
-------

================= ============
**required**      true
**type**          string
**default**       "phantomjs"
================= ============

The browser to run the tests for.

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 3
   :lines: 1-3


steps
-----

================= =================
**required**      false
**type**          [] :doc:`step`
**default**       nil
================= =================

The steps to run for the case.


assertions
----------

================= ====================
**required**      false
**type**          [] :doc:`assertion`
**default**       nil
================= ====================

The assertions to run at the end of the case.
