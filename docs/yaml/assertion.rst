Assertion
=========

An assertion is a validation to run against the current state of an element or page.


text
----

================= ==========
**required**      false
**type**          string
**default**       nil
================= ==========

If provided, asserts that the ``Text`` of the selection matches the provided value.

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 2
   :lines: 33-34


query
-----

================= =============
**required**      false
**type**          [] :doc:`kv`
**default**       nil
================= =============

If provided, validates the url query arguments match each provided query.

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 2-5
   :lines: 38-42


attributes
----------

================= =============
**required**      false
**type**          [] :doc:`kv`
**default**       nil
================= =============

If provided, validates the attributes of the selected elements.

Example
^^^^^^^

.. literalinclude:: example.yml
   :language: yaml
   :emphasize-lines: 2-5
   :lines: 27-31
